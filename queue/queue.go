// Package queue propagates correlation metadata through backend-neutral queue
// metadata maps. New provides Send/Receive; each receipt invokes request identity
// generation anew. Maps are caller-owned and must be stable for reads/exclusively
// owned for writes. No queue backend, authentication, acknowledgment, execution,
// retry, goroutine or cancellation mechanism is introduced.
//
// Deprecated: use github.com/faustbrian/go-correlation/adapters/queue. This
// package remains supported through the documented compatibility interval.
package queue

import (
	"errors"

	correlation "github.com/faustbrian/go-correlation"
)

// ErrInvalidOptions reports invalid queue propagation configuration.
var ErrInvalidOptions = errors.New("queue correlation: invalid options")

// Options configure queue metadata validation and field names.
type Options struct {
	// Codec selects exact metadata names and identifier limits. Empty options
	// choose the root Codec defaults, separately from Factory policy.
	Codec correlation.CodecOptions
}

// Adapter explicitly sends and receives queue metadata. Construct with New;
// shareability depends on the retained Factory generator and caller-owned maps.
type Adapter struct{ propagator *correlation.Propagator }

// New constructs a queue adapter retaining factory. Nil factory returns
// ErrInvalidOptions; codec validation errors are returned without recategorizing.
func New(factory *correlation.Factory, options Options) (*Adapter, error) {
	if factory == nil {
		return nil, ErrInvalidOptions
	}
	codec, err := correlation.NewCodec(options.Codec)
	if err != nil {
		return nil, err
	}
	propagator, _ := correlation.NewPropagator(factory, codec)
	return &Adapter{propagator: propagator}, nil
}

// Send creates a message hop and injects it into an application-owned map.
// Nil metadata/adapter is invalid; populated selected fields refuse overwrite.
// Errors return zero Values but consumed generation state is not restored.
// It does not enqueue or persist a message.
func (adapter *Adapter) Send(metadata map[string]string, parent correlation.Values) (correlation.Values, error) {
	if adapter == nil || metadata == nil {
		return correlation.Values{}, ErrInvalidOptions
	}
	return adapter.propagator.Send(mapCarrier(metadata), parent)
}

// Receive creates a fresh delivery-attempt ID. trusted must be explicitly true
// only after the queue backend's metadata boundary has been authenticated.
// Nil metadata means absence; metadata is never mutated. Malformed untrusted
// fields still fail extraction before trust. The received causation field is not
// reused; accepted inbound request identity becomes the immediate cause.
func (adapter *Adapter) Receive(metadata map[string]string, trusted bool) (correlation.Values, error) {
	if adapter == nil {
		return correlation.Values{}, ErrInvalidOptions
	}
	if metadata == nil {
		metadata = map[string]string{}
	}
	return adapter.propagator.Receive(mapCarrier(metadata), correlation.InboundPolicy{
		TrustCorrelation: trusted, TrustRequestAsCausation: trusted,
	})
}

type mapCarrier map[string]string

func (carrier mapCarrier) Values(key string) []string {
	value, ok := carrier[key]
	if !ok {
		return nil
	}
	return []string{value}
}

func (carrier mapCarrier) Set(key, value string) { carrier[key] = value }
