// Package correlation provides transport-neutral correlation, request, and
// causation identifiers. Identifiers are diagnostic metadata only; they are
// never authentication, authorization, tenancy, replay, or idempotency proof.
// Begin with NewFactory for root/hop values, NewCodec for transport metadata and
// an explicit InboundPolicy for receipt. Named strings and context storage do not
// themselves validate or redact IDs; use Parse helpers at untrusted boundaries.
//
// Values are copied structs containing immutable strings; codecs/strategies have immutable
// configuration. Factories borrow their generator and adapters borrow callbacks
// and mutable carriers, whose synchronization remains caller-owned. The package
// owns no background work, transport, retry, trace provider or shutdown. Calls
// are synchronous with no added deadline or panic recovery for caller code.
package correlation

import (
	"context"
	"errors"
	"fmt"
)

const defaultMaxLength = 128

// ErrInvalidID reports an identifier that violates its validation policy.
var ErrInvalidID = errors.New("correlation: invalid identifier")

// Policy bounds and validates an identifier. The zero value accepts the
// canonical ASCII alphabet [A-Za-z0-9_-] up to 128 bytes.
type Policy struct {
	// MaxLength is the maximum ASCII byte length, not a normalization target.
	// Zero selects 128; positive values through 1024 are supported. Other
	// values reject parsing/construction. Built-in Factory UUID generation is
	// trusted separately and always emits its 36-byte representation.
	MaxLength int
}

// CorrelationID groups work in one logical interaction or workflow. Its empty
// value represents absence in Values, though ParseCorrelationID rejects empty
// input. Direct string conversion bypasses validation and supplies no trust.
//
//revive:disable-next-line:exported semantic name intentionally remains explicit
type CorrelationID string

// RequestID identifies one transport request or delivery attempt. Empty means
// absent metadata; direct conversion bypasses ParseRequestID validation.
type RequestID string

// CausationID identifies the immediate parent request, message, or event, not an
// ancestry chain. Empty means absence; direct conversion bypasses validation.
type CausationID string

// Values carries the three deliberately distinct identifier semantics.
type Values struct {
	// CorrelationID is workflow identity, preserved only by explicit hop/trust
	// policy. Assignment does not validate it.
	CorrelationID CorrelationID
	// RequestID is this hop/attempt identity, newly generated at each boundary.
	RequestID RequestID
	// CausationID is the immediate parent's request identity, or empty at a root.
	CausationID CausationID
}

// ParseCorrelationID accepts nonempty canonical ASCII text unchanged under
// policy. Invalid alphabet, length or policy returns an empty ID and an error
// wrapping ErrInvalidID, without rejected input. It does not assign trust.
func ParseCorrelationID(value string, policy Policy) (CorrelationID, error) {
	if err := validate(value, policy); err != nil {
		return "", err
	}
	return CorrelationID(value), nil
}

// ParseRequestID validates nonempty ASCII text exactly as ParseCorrelationID,
// returning an empty ID and ErrInvalidID-wrapping error on failure.
func ParseRequestID(value string, policy Policy) (RequestID, error) {
	if err := validate(value, policy); err != nil {
		return "", err
	}
	return RequestID(value), nil
}

// ParseCausationID validates nonempty ASCII text exactly as ParseCorrelationID,
// returning an empty ID and ErrInvalidID-wrapping error on failure.
func ParseCausationID(value string, policy Policy) (CausationID, error) {
	if err := validate(value, policy); err != nil {
		return "", err
	}
	return CausationID(value), nil
}

// MustCorrelationID is ParseCorrelationID for static configuration and tests.
// It panics with the validation error for invalid input or policy.
func MustCorrelationID(value string, policy Policy) CorrelationID {
	id, err := ParseCorrelationID(value, policy)
	if err != nil {
		panic(err)
	}
	return id
}

// MustRequestID is ParseRequestID for static configuration and tests.
// It panics with the validation error for invalid input or policy.
func MustRequestID(value string, policy Policy) RequestID {
	id, err := ParseRequestID(value, policy)
	if err != nil {
		panic(err)
	}
	return id
}

// MustCausationID is ParseCausationID for static configuration and tests.
// It panics with the validation error for invalid input or policy.
func MustCausationID(value string, policy Policy) CausationID {
	id, err := ParseCausationID(value, policy)
	if err != nil {
		panic(err)
	}
	return id
}

// String returns the underlying text without validation or redaction.
func (id CorrelationID) String() string { return string(id) }

// String returns the underlying text without validation or redaction.
func (id RequestID) String() string { return string(id) }

// String returns the underlying text without validation or redaction.
func (id CausationID) String() string { return string(id) }

type contextKey struct{}

// WithValues returns a derived context carrying a copy of values under a
// package-private key, replacing prior Values at that key without merging or
// validation. It preserves parent cancellation/deadlines and panics for nil ctx.
func WithValues(ctx context.Context, values Values) context.Context {
	return context.WithValue(ctx, contextKey{}, values)
}

// FromContext returns a copy of Values installed by WithValues. Nil or absent
// context returns zero Values, false; explicitly stored zero Values returns true.
// It does not validate, authenticate or redact the stored strings.
func FromContext(ctx context.Context) (Values, bool) {
	if ctx == nil {
		return Values{}, false
	}
	values, ok := ctx.Value(contextKey{}).(Values)
	return values, ok
}

func validate(value string, policy Policy) error {
	maximum := policy.MaxLength
	if maximum == 0 {
		maximum = defaultMaxLength
	}
	if maximum > 1024 {
		return fmt.Errorf("%w: length", ErrInvalidID)
	}
	if len(value) == 0 {
		return fmt.Errorf("%w: length", ErrInvalidID)
	}
	if len(value) > maximum {
		return fmt.Errorf("%w: length", ErrInvalidID)
	}
	for index := range len(value) {
		char := value[index]
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') && char != '-' && char != '_' {
			return fmt.Errorf("%w: alphabet", ErrInvalidID)
		}
	}
	return nil
}
