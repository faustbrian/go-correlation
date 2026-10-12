package correlation

import (
	"errors"
	"fmt"
)

// ErrInvalidPropagator reports missing explicit propagation dependencies.
var ErrInvalidPropagator = errors.New("correlation: invalid propagator")

// Propagator explicitly composes hop generation with a transport codec. It
// retains both dependencies, owns no carrier or transport resources and inherits
// Factory's generator synchronization and Carrier's caller-owned mutation rules.
type Propagator struct {
	factory *Factory
	codec   *Codec
}

// NewPropagator rejects nil dependencies with ErrInvalidPropagator. Supply
// constructed Factory/Codec values; it does not validate their internal state.
func NewPropagator(factory *Factory, codec *Codec) (*Propagator, error) {
	if factory == nil || codec == nil {
		return nil, ErrInvalidPropagator
	}
	return &Propagator{factory: factory, codec: codec}, nil
}

// Send creates the next hop before injecting its values. Failure returns zero
// Values and may have consumed generator state; injection errors are wrapped.
// It performs no send, retry or acknowledgment, and cannot roll back callbacks.
func (propagator *Propagator) Send(carrier Carrier, parent Values) (Values, error) {
	if propagator == nil || propagator.factory == nil || propagator.codec == nil {
		return Values{}, ErrInvalidPropagator
	}
	values, err := propagator.factory.Next(parent)
	if err != nil {
		return Values{}, err
	}
	if err := propagator.codec.Inject(carrier, values); err != nil {
		return Values{}, fmt.Errorf("send correlation: %w", err)
	}
	return values, nil
}

// Receive extracts untrusted metadata and applies an explicit trust policy
// while creating a fresh delivery-attempt request ID. Extraction checks every
// carried field before trust, so malformed untrusted metadata still fails here.
// Failures return zero Values with the extraction/generation cause preserved.
func (propagator *Propagator) Receive(carrier Carrier, policy InboundPolicy) (Values, error) {
	if propagator == nil || propagator.factory == nil || propagator.codec == nil {
		return Values{}, ErrInvalidPropagator
	}
	inbound, err := propagator.codec.Extract(carrier)
	if err != nil {
		return Values{}, fmt.Errorf("receive correlation: %w", err)
	}
	return propagator.factory.Accept(inbound, policy)
}
