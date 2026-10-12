package correlation_test

import (
	"testing"

	correlation "github.com/faustbrian/go-correlation"
)

// FuzzParseCorrelationID supplies arbitrary text under a 128-byte ASCII policy.
// Successfully parsed text must be unchanged; rejection is not a validity oracle
// beyond this package's policy, and arbitrary named-string casts are not tested.
func FuzzParseCorrelationID(fuzz *testing.F) {
	for _, seed := range []string{"flow", "request-123", "bad value", "control\n", "", "d1_stable"} {
		fuzz.Add(seed)
	}
	fuzz.Fuzz(func(t *testing.T, value string) {
		id, err := correlation.ParseCorrelationID(value, correlation.Policy{MaxLength: 128})
		if err == nil && id.String() != value {
			t.Fatalf("canonical text changed from %q to %q", value, id)
		}
	})
}

// FuzzCarrierExtraction mutates three fields and at most sixteen added duplicate
// strings. Successful extraction must reinject into a fresh carrier without error;
// input allocation is outside the bounded eight-value Codec acceptance rule.
func FuzzCarrierExtraction(fuzz *testing.F) {
	fuzz.Add("flow", "request", "parent", "flow", uint8(1))
	fuzz.Add("bad value", "", "control\n", "conflict", uint8(9))
	fuzz.Fuzz(func(t *testing.T, correlationValue, requestValue, causationValue, duplicateValue string, count uint8) {
		codec, err := correlation.NewCodec(correlation.CodecOptions{Policy: correlation.Policy{MaxLength: 128}})
		if err != nil {
			t.Fatal(err)
		}
		carrier := memoryCarrier{}
		if correlationValue != "" {
			carrier[correlation.DefaultCorrelationField] = []string{correlationValue}
		}
		if requestValue != "" {
			carrier[correlation.DefaultRequestField] = []string{requestValue}
		}
		if causationValue != "" {
			carrier[correlation.DefaultCausationField] = []string{causationValue}
		}
		for range min(int(count), 16) {
			carrier[correlation.DefaultCorrelationField] = append(
				carrier[correlation.DefaultCorrelationField], duplicateValue,
			)
		}
		values, err := codec.Extract(carrier)
		if err == nil {
			output := memoryCarrier{}
			if err := codec.Inject(output, values); err != nil {
				t.Fatalf("valid extracted values failed reinjection: %v", err)
			}
		}
	})
}
