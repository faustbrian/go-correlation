package correlation

import (
	"bufio"
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"io"

	identifieruuid "github.com/faustbrian/go-identifier/v2/uuid"
)

const defaultEntropyBufferSize = 4 << 10

var (
	// ErrInvalidFactory reports invalid factory configuration.
	ErrInvalidFactory = errors.New("correlation: invalid factory")
	// ErrGeneration reports a generator failure or invalid generated value.
	ErrGeneration = errors.New("correlation: generation failed")
)

// Generator supplies canonical random identifier text. It is structurally
// compatible with identifier.Generator[string] from identifier. Factory borrows
// it: implementations own synchronization, randomness/uniqueness and lifetime.
// Calls are synchronous, may consume state on failure and must not assume retry.
type Generator interface {
	// New returns nonempty text accepted by Factory's Policy or an error.
	// Factories validate successful custom results, wrap errors without
	// redaction, and neither recover panics nor impose cancellation.
	New() (string, error)
}

// GeneratorFunc adapts a borrowed synchronous function to Generator. Nil
// functions panic when called directly and are rejected by NewFactory.
type GeneratorFunc func() (string, error)

// New calls the wrapped function, preserving its result, error and panic.
// The function must supply its own synchronization when shared.
func (function GeneratorFunc) New() (string, error) { return function() }

// FactoryOptions configure an immutable Factory.
type FactoryOptions struct {
	// Policy validates parents, trusted inbound values and custom generated
	// results; it does not truncate/revalidate built-in 36-byte UUID results.
	Policy Policy
	// Generator is borrowed for the Factory lifetime, never copied or closed.
	// Literal nil selects the cryptographic default; typed nil is invalid.
	Generator Generator
}

// InboundPolicy explicitly identifies which inbound semantics cross a trust
// boundary. Request IDs are never preserved because each hop gets a new one.
type InboundPolicy struct {
	// TrustCorrelation permits retaining a validated nonempty inbound workflow
	// ID. It records the caller's boundary decision, not authentication.
	TrustCorrelation bool
	// TrustRequestAsCausation permits a validated inbound request ID as the
	// immediate cause. Inbound causation itself is always ignored.
	TrustRequestAsCausation bool
}

// Factory creates hop identifiers without global mutable state. Use NewFactory;
// its configuration is immutable, but concurrent use requires a concurrency-safe
// Generator (the built-in UUID generator is safe). No uniqueness ledger or
// collision detection is owned; custom generators determine those properties.
type Factory struct {
	policy            Policy
	generator         Generator
	validateGenerated bool
}

// NewFactory validates policy, then constructs a factory. The default is a
// buffered crypto/rand-backed UUIDv4 generator from Identifier v2.0.1. Its
// canonical 36-byte output is trusted, even with Policy.MaxLength below 36.
// Custom generators are retained and their successful output is validated.
// Invalid configuration returns nil with ErrInvalidFactory; no generator is run.
func NewFactory(options FactoryOptions) (*Factory, error) {
	if err := validatePolicy(options.Policy); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidFactory, err)
	}
	validateGenerated := options.Generator != nil
	if validateGenerated && nilLike(options.Generator) {
		return nil, fmt.Errorf("%w: nil generator", ErrInvalidFactory)
	}
	if options.Generator == nil {
		options.Generator = newBufferedUUIDGenerator(cryptorand.Reader)
	}
	return &Factory{
		policy:            options.Policy,
		generator:         options.Generator,
		validateGenerated: validateGenerated,
	}, nil
}

type uuidGenerator struct{ generator *identifieruuid.V4Generator }

// Identifier's generator owns synchronization around this borrowed buffer. Keeping
// the 4 KiB buffer per factory amortizes entropy reads without ambient state;
// TestDefaultGeneratorAmortizesEntropyReads protects that boundary.
func newBufferedUUIDGenerator(reader io.Reader) *uuidGenerator {
	return &uuidGenerator{generator: identifieruuid.NewV4Generator(
		bufio.NewReaderSize(reader, defaultEntropyBufferSize),
	)}
}

func (generator *uuidGenerator) New() (string, error) {
	id, err := generator.generator.New()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// Create generates correlation first, then request identity, with no causation.
// It returns zero Values on any ErrGeneration-wrapping failure, retaining custom
// generator causes; already consumed generator state is not rolled back. Nil or
// uninitialized factories return ErrGeneration here.
func (factory *Factory) Create() (Values, error) {
	correlationID, err := factory.newCorrelationID()
	if err != nil {
		return Values{}, err
	}
	requestID, err := factory.newRequestID()
	if err != nil {
		return Values{}, err
	}
	return Values{CorrelationID: correlationID, RequestID: requestID}, nil
}

// Start creates a new correlation and request identifier.
//
// Deprecated: use Create.
func (factory *Factory) Start() (Values, error) {
	return factory.Create()
}

// Next preserves correlation, creates a request ID, and makes the prior
// request the immediate cause. Parent causation is ignored. Parent validation
// precedes generation; failures return zero Values wrapping ErrGeneration and
// the validation/generator cause. The receiver must be nonnil.
func (factory *Factory) Next(parent Values) (Values, error) {
	correlationID, err := ParseCorrelationID(parent.CorrelationID.String(), factory.policy)
	if err != nil {
		return Values{}, fmt.Errorf("%w: parent correlation: %w", ErrGeneration, err)
	}
	causationID, err := ParseCausationID(parent.RequestID.String(), factory.policy)
	if err != nil {
		return Values{}, fmt.Errorf("%w: parent request: %w", ErrGeneration, err)
	}
	requestID, err := factory.newRequestID()
	if err != nil {
		return Values{}, err
	}
	return Values{
		CorrelationID: correlationID,
		RequestID:     requestID,
		CausationID:   causationID,
	}, nil
}

// Accept starts a receiving hop. No inbound value is used unless its exact
// semantic trust is enabled. Missing trusted values still select generation;
// inbound causation is ignored. Successful receipt generates a request ID; errors return
// zero Values on failure without rolling back consumed generator state. Use a
// constructed nonnil factory: a nil receiver panics when TrustCorrelation is true
// and inbound.CorrelationID is nonempty; other nil paths fail generation.
func (factory *Factory) Accept(inbound Values, policy InboundPolicy) (Values, error) {
	var correlationID CorrelationID
	var err error
	if policy.TrustCorrelation && inbound.CorrelationID != "" {
		correlationID, err = ParseCorrelationID(inbound.CorrelationID.String(), factory.policy)
		if err != nil {
			return Values{}, fmt.Errorf("%w: inbound correlation: %w", ErrGeneration, err)
		}
	} else {
		correlationID, err = factory.newCorrelationID()
		if err != nil {
			return Values{}, err
		}
	}

	requestID, err := factory.newRequestID()
	if err != nil {
		return Values{}, err
	}
	values := Values{CorrelationID: correlationID, RequestID: requestID}
	if policy.TrustRequestAsCausation && inbound.RequestID != "" {
		values.CausationID, err = ParseCausationID(inbound.RequestID.String(), factory.policy)
		if err != nil {
			return Values{}, fmt.Errorf("%w: inbound request: %w", ErrGeneration, err)
		}
	}
	return values, nil
}

func (factory *Factory) newCorrelationID() (CorrelationID, error) {
	value, err := factory.generate()
	if err != nil {
		return "", err
	}
	return CorrelationID(value), nil
}

func (factory *Factory) newRequestID() (RequestID, error) {
	value, err := factory.generate()
	if err != nil {
		return "", err
	}
	return RequestID(value), nil
}

func (factory *Factory) generate() (string, error) {
	if factory == nil || factory.generator == nil {
		return "", fmt.Errorf("%w: nil factory", ErrGeneration)
	}
	value, err := factory.generator.New()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrGeneration, err)
	}
	// The default UUID path is a trusted canonical producer. This deliberately
	// differs from custom output validation, including policy length checks.
	if factory.validateGenerated {
		if err := validate(value, factory.policy); err != nil {
			return "", fmt.Errorf("%w: %w", ErrGeneration, err)
		}
	}
	return value, nil
}

func validatePolicy(policy Policy) error {
	maximum := policy.MaxLength
	if maximum == 0 {
		maximum = defaultMaxLength
	}
	if maximum < 1 || maximum > 1024 {
		return fmt.Errorf("maximum length must be between 1 and 1024")
	}
	return nil
}
