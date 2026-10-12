// Package httpcorrelation explicitly propagates correlation metadata over
// HTTP. New composes an explicit Factory with Wrap for inbound handlers and
// Inject for outbound hops. Middleware installation, peer authentication,
// transport execution, cancellation and retry remain application-owned.
// Configuration may be shared when the generator/Trust callback are safe to
// share; request/response header maps need exclusive mutation ownership.
// No global middleware, goroutine, body consumption or shutdown is introduced.
//
// Deprecated: use github.com/faustbrian/go-correlation/adapters/http. This
// package remains supported through the documented compatibility interval.
package httpcorrelation

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	correlation "github.com/faustbrian/go-correlation"
)

const (
	// CorrelationHeader carries the logical workflow identifier.
	CorrelationHeader = "X-Correlation-ID"
	// RequestHeader carries the current HTTP hop identifier.
	RequestHeader = "X-Request-ID"
	// CausationHeader carries the immediate parent identifier.
	CausationHeader   = "X-Causation-ID"
	correlationMapKey = "X-Correlation-Id"
	requestMapKey     = "X-Request-Id"
	causationMapKey   = "X-Causation-Id"
	maxHeaderValues   = 8
)

// InvalidPolicy controls malformed or ambiguous inbound metadata.
type InvalidPolicy uint8

const (
	// ReplaceInvalid discards all metadata on extraction failure and creates a
	// fresh root. With no Trust callback, inbound parsing is skipped entirely.
	ReplaceInvalid InvalidPolicy = iota
	// RejectInvalid returns HTTP 400 before application code runs.
	RejectInvalid
)

// ErrInvalidOptions reports invalid HTTP adapter configuration.
var ErrInvalidOptions = errors.New("http correlation: invalid options")

// Options configure immutable HTTP propagation policy.
type Options struct {
	// Policy configures inbound/outbound codec validation, independently of the
	// supplied Factory's policy; zero selects correlation's defaults.
	Policy correlation.Policy
	// Invalid selects replacement (zero) or HTTP 400 for malformed metadata.
	Invalid InvalidPolicy
	// Trust is a borrowed synchronous peer-authentication decision. Nil trusts
	// no inbound IDs. It may run concurrently and can block or panic; the
	// adapter adds no recovery, deadline or authentication of its own.
	Trust func(*http.Request) bool
}

// Middleware owns explicit HTTP extraction, trust, context, and injection.
// Use New; immutable configuration is shareable subject to Factory/Trust safety.
// It owns neither request/response lifetime nor mutable header synchronization.
type Middleware struct {
	factory    *correlation.Factory
	codec      *correlation.Codec
	propagator *correlation.Propagator
	options    Options
}

// New constructs an HTTP adapter without installing global middleware. Nil
// factory, unknown Invalid mode or invalid codec policy returns an error; factory
// and Trust are borrowed. Supply a constructed Factory, not its zero value.
func New(factory *correlation.Factory, options Options) (*Middleware, error) {
	if factory == nil || options.Invalid > RejectInvalid {
		return nil, ErrInvalidOptions
	}
	codec, err := correlation.NewCodec(correlation.CodecOptions{
		Policy: options.Policy, CorrelationField: CorrelationHeader,
		RequestField: RequestHeader, CausationField: CausationHeader,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidOptions, err)
	}
	propagator, _ := correlation.NewPropagator(factory, codec)
	return &Middleware{factory: factory, codec: codec, propagator: propagator, options: options}, nil
}

// Wrap creates one fresh request ID per HTTP request. Valid inbound metadata
// is preserved only when Trust explicitly accepts the immediate peer. Extraction
// failure under RejectInvalid returns 400 before next; generation failure returns
// 500. In both successful and rejected flows, canonical request/response header
// keys are written; this is not removal of every alternate-cased map key.
// Request headers must be writable/nonnil and exclusively owned for the call.
// Wrap derives context but does not deep-copy the request, consume its body,
// cancel work or recover callback/handler panics. Nil next yields a 500 handler.
func (middleware *Middleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if middleware == nil || next == nil || request == nil {
			http.Error(writer, "internal server error", http.StatusInternalServerError)
			return
		}
		values, rejected, err := middleware.receive(request)
		if err != nil {
			http.Error(writer, "internal server error", http.StatusInternalServerError)
			return
		}
		if rejected {
			setHeaders(request.Header, values)
			setHeaders(writer.Header(), values)
			http.Error(writer, "invalid correlation metadata", http.StatusBadRequest)
			return
		}

		setHeaders(request.Header, values)
		setHeaders(writer.Header(), values)
		ctx := correlation.WithValues(request.Context(), values)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (middleware *Middleware) receive(
	request *http.Request,
) (correlation.Values, bool, error) {
	if middleware.options.Invalid == ReplaceInvalid && middleware.options.Trust == nil {
		values, err := middleware.factory.Create()

		return values, false, err
	}
	inbound, err := middleware.codec.Extract(headerCarrier{request.Header})
	if err != nil && middleware.options.Invalid == RejectInvalid {
		values, generationErr := middleware.factory.Create()
		if generationErr != nil {
			return correlation.Values{}, false, generationErr
		}

		return values, true, nil
	}
	if err != nil {
		inbound = correlation.Values{}
	}
	trusted := middleware.options.Trust != nil && middleware.options.Trust(request)

	values, err := middleware.factory.Accept(inbound, correlation.InboundPolicy{
		TrustCorrelation: trusted, TrustRequestAsCausation: trusted,
	})

	return values, false, err
}

// Inject creates and injects the next outbound HTTP hop without sending it.
// It mutates request.Header (allocating it if nil), preserving unrelated fields
// and refusing populated correlation fields through Codec. Failures return zero
// Values but may allocate headers/consume generator state; no rollback occurs.
// Nil request or Middleware returns ErrInvalidOptions.
func (middleware *Middleware) Inject(request *http.Request, parent correlation.Values) (correlation.Values, error) {
	if middleware == nil || request == nil {
		return correlation.Values{}, ErrInvalidOptions
	}
	if request.Header == nil {
		request.Header = make(http.Header)
	}
	return middleware.propagator.Send(headerCarrier{request.Header}, parent)
}

type headerCarrier struct{ header http.Header }

// Direct Header maps may contain noncanonical casing. Scan all case-equivalent
// keys so conflicts cannot disappear behind Header.Get's canonical lookup; keep
// one excess value to let Codec reject before collecting an unbounded slice.
func (carrier headerCarrier) Values(key string) []string {
	var values []string
	for name, entries := range carrier.header {
		if strings.EqualFold(name, key) {
			for _, entry := range entries {
				values = append(values, entry)
				if len(values) > maxHeaderValues {
					return values
				}
			}
		}
	}
	return values
}

func (carrier headerCarrier) Set(key, value string) { carrier.header.Set(key, value) }

func setHeaders(header http.Header, values correlation.Values) {
	header[correlationMapKey] = []string{values.CorrelationID.String()}
	header[requestMapKey] = []string{values.RequestID.String()}
	if values.CausationID != "" {
		header[causationMapKey] = []string{values.CausationID.String()}
	} else {
		delete(header, causationMapKey)
	}
}
