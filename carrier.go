package correlation

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	// DefaultCorrelationField is the transport-neutral correlation key.
	DefaultCorrelationField = "correlation_id"
	// DefaultRequestField is the transport-neutral request key.
	DefaultRequestField = "request_id"
	// DefaultCausationField is the transport-neutral causation key.
	DefaultCausationField = "causation_id"
)

var (
	// ErrInvalidCarrier reports malformed, unbounded, or unsupported metadata.
	ErrInvalidCarrier = errors.New("correlation: invalid carrier")
	// ErrConflictingCarrier reports more than one distinct value for a field.
	ErrConflictingCarrier = errors.New("correlation: conflicting carrier values")
	// ErrCarrierOverwrite reports injection into an already populated field.
	ErrCarrierOverwrite = errors.New("correlation: carrier overwrite")
)

// Carrier is an explicit transport metadata boundary. Values must return a
// copy or immutable view stable for the call. Set replaces one field with exactly
// one value. The caller owns carrier lifetime, synchronization, name matching
// and callback behavior; Codec provides no transaction, timeout or panic recovery.
type Carrier interface {
	// Values returns all values for key, preserving duplicates; absence is a
	// zero-length slice, distinct from a present empty string. Providers should
	// bound collection before allocation for hostile transport metadata.
	Values(key string) []string
	// Set replaces key with one value. It must not change unrelated fields;
	// injection preflights before invoking it but cannot undo provider failures.
	Set(key, value string)
}

// CodecOptions configure immutable carrier field names and validation.
type CodecOptions struct {
	// Policy validates carried identifier text; zero selects default limits.
	Policy Policy
	// CorrelationField selects a distinct ASCII name of at most 64 bytes;
	// empty selects DefaultCorrelationField.
	CorrelationField string
	// RequestField is a distinct ASCII name; empty selects DefaultRequestField.
	RequestField string
	// CausationField is a distinct ASCII name; empty selects DefaultCausationField.
	CausationField string
}

// Codec injects and extracts typed values without assigning trust. Construct
// with NewCodec; immutable configuration can be shared, but carriers must be
// synchronized independently. Zero Codec values are not configured codecs.
type Codec struct {
	policy           Policy
	correlationField string
	requestField     string
	causationField   string
}

// NewCodec validates and copies carrier configuration, returning nil and an
// ErrInvalidCarrier-wrapping error on invalid policy or duplicate/invalid names.
// Field uniqueness is exact, not case-folded; custom names must also be distinct
// under the selected transport's own name comparison.
func NewCodec(options CodecOptions) (*Codec, error) {
	if err := validatePolicy(options.Policy); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidCarrier, err)
	}
	if options.CorrelationField == "" {
		options.CorrelationField = DefaultCorrelationField
	}
	if options.RequestField == "" {
		options.RequestField = DefaultRequestField
	}
	if options.CausationField == "" {
		options.CausationField = DefaultCausationField
	}
	fields := []string{options.CorrelationField, options.RequestField, options.CausationField}
	seen := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		if err := validate(field, Policy{MaxLength: 64}); err != nil {
			return nil, fmt.Errorf("%w: field name", ErrInvalidCarrier)
		}
		if _, exists := seen[field]; exists {
			return nil, fmt.Errorf("%w: duplicate field name", ErrInvalidCarrier)
		}
		seen[field] = struct{}{}
	}
	return &Codec{
		policy: options.Policy, correlationField: options.CorrelationField,
		requestField: options.RequestField, causationField: options.CausationField,
	}, nil
}

// Extract parses values but deliberately does not decide whether to trust
// them. Pass the result to Factory.Accept with an explicit InboundPolicy.
// Absent fields are allowed; at most eight repeated values must all be equal and
// nonempty. Invalid/conflicting metadata returns zero Values, not a partial parse.
// Nil codec or nil/typed-nil carrier returns ErrInvalidCarrier. No writes occur.
func (codec *Codec) Extract(carrier Carrier) (Values, error) {
	if codec == nil || nilLike(carrier) {
		return Values{}, fmt.Errorf("%w: nil codec or carrier", ErrInvalidCarrier)
	}
	correlationText, err := codec.one(carrier, codec.correlationField)
	if err != nil {
		return Values{}, err
	}
	requestText, err := codec.one(carrier, codec.requestField)
	if err != nil {
		return Values{}, err
	}
	causationText, err := codec.one(carrier, codec.causationField)
	if err != nil {
		return Values{}, err
	}

	var values Values
	if correlationText != "" {
		values.CorrelationID, err = ParseCorrelationID(correlationText, codec.policy)
	}
	if err == nil && requestText != "" {
		values.RequestID, err = ParseRequestID(requestText, codec.policy)
	}
	if err == nil && causationText != "" {
		values.CausationID, err = ParseCausationID(causationText, codec.policy)
	}
	if err != nil {
		return Values{}, fmt.Errorf("%w: %w", ErrInvalidCarrier, err)
	}
	return values, nil
}

// Inject installs non-empty values and refuses to overwrite any populated
// field, even when the existing value is malformed. Empty output fields are
// skipped, including overwrite checks. All nonempty values/overwrites are checked
// before the first Set; this is not atomic against concurrent carrier mutation
// or a provider panic. Validation/overwrite failure makes no Set calls.
func (codec *Codec) Inject(carrier Carrier, values Values) error {
	if codec == nil || nilLike(carrier) {
		return fmt.Errorf("%w: nil codec or carrier", ErrInvalidCarrier)
	}
	fields := []struct {
		name  string
		value string
	}{
		{codec.correlationField, values.CorrelationID.String()},
		{codec.requestField, values.RequestID.String()},
		{codec.causationField, values.CausationID.String()},
	}
	for _, field := range fields {
		if field.value == "" {
			continue
		}
		if err := validate(field.value, codec.policy); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidCarrier, err)
		}
		if len(carrier.Values(field.name)) != 0 {
			return fmt.Errorf("%w: %s", ErrCarrierOverwrite, field.name)
		}
	}
	// Separate preflight prevents a later invalid/occupied field from leaving
	// earlier fields written. Carrier callbacks themselves are not transactional.
	for _, field := range fields {
		if field.value != "" {
			carrier.Set(field.name, field.value)
		}
	}
	return nil
}

func (codec *Codec) one(carrier Carrier, field string) (string, error) {
	values := carrier.Values(field)
	if len(values) == 0 {
		return "", nil
	}
	if len(values) > 8 {
		return "", fmt.Errorf("%w: too many %s values", ErrInvalidCarrier, field)
	}
	first := values[0]
	if first == "" {
		return "", fmt.Errorf("%w: empty %s", ErrInvalidCarrier, field)
	}
	for _, value := range values[1:] {
		if value != first {
			return "", fmt.Errorf("%w: %s", ErrConflictingCarrier, field)
		}
	}
	return first, nil
}

// Interface presence alone does not distinguish a usable callback/carrier from
// typed-nil pointers, maps, slices, channels or functions (lifecycle_test.go).
func nilLike(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	//exhaustive:ignore only nil-capable kinds need special handling
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}
