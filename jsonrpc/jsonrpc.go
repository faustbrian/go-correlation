// Package jsonrpc propagates correlation through an explicit JSON-RPC
// metadata object. New composes Send/Receive over caller-owned Metadata; it never
// edits protocol envelopes or JSON-RPC request IDs. Callers decode envelopes,
// preserve duplicate metadata members, authenticate trust and synchronize maps.
// No transport, background work, retry or cancellation is introduced.
//
// Deprecated: use github.com/faustbrian/go-correlation/adapters/jsonrpc. This
// package remains supported through the documented compatibility interval.
package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"

	correlation "github.com/faustbrian/go-correlation"
)

const (
	// CorrelationField carries the logical workflow identifier.
	CorrelationField = correlation.DefaultCorrelationField
	// RequestField carries the current JSON-RPC hop identifier.
	RequestField = correlation.DefaultRequestField
	// CausationField carries the immediate parent identifier.
	CausationField    = correlation.DefaultCausationField
	maxMetadataValues = 8
	maxEncodedIDBytes = 1026
)

var (
	// ErrInvalidOptions reports incomplete adapter configuration.
	ErrInvalidOptions = errors.New("jsonrpc correlation: invalid options")
	// ErrMalformedMetadata reports non-string or oversized JSON metadata.
	ErrMalformedMetadata = errors.New("jsonrpc correlation: malformed metadata")
)

// Metadata preserves repeated JSON members supplied by a strict envelope
// decoder so conflicting duplicates cannot be silently collapsed. Assignment
// aliases the map and RawMessage byte slices. Keep them stable during receipt;
// Send replaces only selected keys. Whole-envelope allocation is caller-bounded.
type Metadata map[string][]json.RawMessage

// Options configure JSON-RPC metadata validation and names.
type Options struct {
	// Codec selects root validation and distinct exact metadata names; empty
	// names select this package's default fields.
	Codec correlation.CodecOptions
}

// Adapter explicitly sends and receives JSON-RPC metadata. Use New; configuration
// is immutable but concurrent sharing inherits Factory generator safety and
// requires independently synchronized metadata maps/bytes.
type Adapter struct {
	propagator *correlation.Propagator
	fields     [3]string
}

// New constructs a JSON-RPC metadata adapter, retaining the supplied Factory.
// Nil factory returns ErrInvalidOptions; invalid codec options preserve Codec's
// error. No metadata is read and no generator is called at construction.
func New(factory *correlation.Factory, options Options) (*Adapter, error) {
	if factory == nil {
		return nil, ErrInvalidOptions
	}
	codec, err := correlation.NewCodec(options.Codec)
	if err != nil {
		return nil, err
	}
	if options.Codec.CorrelationField == "" {
		options.Codec.CorrelationField = CorrelationField
	}
	if options.Codec.RequestField == "" {
		options.Codec.RequestField = RequestField
	}
	if options.Codec.CausationField == "" {
		options.Codec.CausationField = CausationField
	}
	propagator, _ := correlation.NewPropagator(factory, codec)
	return &Adapter{
		propagator: propagator,
		fields: [3]string{
			options.Codec.CorrelationField,
			options.Codec.RequestField,
			options.Codec.CausationField,
		},
	}, nil
}

// Send creates and injects a child JSON-RPC metadata hop, replacing only selected
// absent keys with JSON strings. Nil metadata/adapter returns ErrInvalidOptions;
// other failures return zero Values and may already consume generator state.
// Metadata mutation follows Codec preflight, not a transport transaction.
func (adapter *Adapter) Send(metadata Metadata, parent correlation.Values) (correlation.Values, error) {
	if adapter == nil || metadata == nil {
		return correlation.Values{}, ErrInvalidOptions
	}
	return adapter.propagator.Send(rawCarrier(metadata), parent)
}

// Receive validates metadata and creates a fresh receiving request ID. Nil
// metadata means absent fields. Each selected field allows at most eight raw
// values, each at most 1026 encoded bytes including quotes/escapes/whitespace;
// whole maps and unrelated members are not bounded by this check. Invalid JSON
// or incompatible types return ErrMalformedMetadata; decoded empty text (including
// null) fails Codec. Trust enables only correlation and request-as-causation after
// all selected metadata is validated. No mutation or protocol-ID mapping occurs.
func (adapter *Adapter) Receive(metadata Metadata, trusted bool) (correlation.Values, error) {
	if adapter == nil {
		return correlation.Values{}, ErrInvalidOptions
	}
	if metadata == nil {
		metadata = Metadata{}
	}
	// Check count/encoded length before decoding to bound this adapter's work.
	// Envelope decoding already happened outside it; unknown members stay opaque.
	for _, field := range adapter.fields {
		if len(metadata[field]) > maxMetadataValues {
			return correlation.Values{}, fmt.Errorf("%w: too many %s values", correlation.ErrInvalidCarrier, field)
		}
		for _, raw := range metadata[field] {
			var value string
			if len(raw) > maxEncodedIDBytes {
				return correlation.Values{}, fmt.Errorf("%w: %s", ErrMalformedMetadata, field)
			}
			if json.Unmarshal(raw, &value) != nil {
				return correlation.Values{}, fmt.Errorf("%w: %s", ErrMalformedMetadata, field)
			}
		}
	}
	return adapter.propagator.Receive(rawCarrier(metadata), correlation.InboundPolicy{
		TrustCorrelation: trusted, TrustRequestAsCausation: trusted,
	})
}

type rawCarrier Metadata

func (carrier rawCarrier) Values(key string) []string {
	rawValues := carrier[key]
	if len(rawValues) > maxMetadataValues {
		rawValues = rawValues[:maxMetadataValues+1]
	}
	values := make([]string, 0, len(rawValues))
	for _, raw := range rawValues {
		var value string
		_ = json.Unmarshal(raw, &value)
		values = append(values, value)
	}
	return values
}

func (carrier rawCarrier) Set(key, value string) {
	raw, _ := json.Marshal(value)
	carrier[key] = []json.RawMessage{raw}
}
