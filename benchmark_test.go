package correlation_test

import (
	"strconv"
	"strings"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
)

type benchmarkGenerator uint64

func (generator *benchmarkGenerator) New() (string, error) {
	*generator++
	return "id_" + strconv.FormatUint(uint64(*generator), 10), nil
}

// BenchmarkCarrierRejectsOversizedValue times rejecting one 4096-byte value
// against default 128-byte policy. b.Loop excludes codec/carrier construction;
// this measures rejection, not parsing or allocation of an external envelope.
func BenchmarkCarrierRejectsOversizedValue(b *testing.B) {
	codec, _ := correlation.NewCodec(correlation.CodecOptions{})
	carrier := memoryCarrier{
		correlation.DefaultCorrelationField: {strings.Repeat("a", 4096)},
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := codec.Extract(carrier); err == nil {
			b.Fatal("oversized carrier was accepted")
		}
	}
}

// BenchmarkFactoryNext times parent validation and one custom sequence-generated
// child ID. b.Loop excludes configuration; its cheap mutable generator is not an
// equivalent workload to crypto/rand-backed default generation.
func BenchmarkFactoryNext(b *testing.B) {
	generator := new(benchmarkGenerator)
	factory, _ := correlation.NewFactory(correlation.FactoryOptions{Generator: generator})
	parent := correlation.Values{CorrelationID: "flow", RequestID: "parent"}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := factory.Next(parent); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDefaultFactoryCreate times two buffered UUIDv4 generations per root.
// b.Loop excludes Factory construction; entropy-buffer refill remains in the loop.
func BenchmarkDefaultFactoryCreate(b *testing.B) {
	factory, _ := correlation.NewFactory(correlation.FactoryOptions{})
	b.ReportAllocs()
	for b.Loop() {
		if _, err := factory.Create(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDefaultFactoryNext times parent validation and one buffered UUIDv4
// generation. b.Loop excludes configuration, not entropy-buffer refills.
func BenchmarkDefaultFactoryNext(b *testing.B) {
	factory, _ := correlation.NewFactory(correlation.FactoryOptions{})
	parent := correlation.Values{CorrelationID: "flow", RequestID: "parent"}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := factory.Next(parent); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCarrierRoundTrip times a fresh in-memory carrier, three-value injection
// and extraction. b.Loop excludes codec/value setup, not per-iteration map/slice
// allocation; it does not include transport encoding, I/O or authentication.
func BenchmarkCarrierRoundTrip(b *testing.B) {
	codec, _ := correlation.NewCodec(correlation.CodecOptions{})
	values := correlation.Values{CorrelationID: "flow", RequestID: "request", CausationID: "parent"}
	b.ReportAllocs()
	for b.Loop() {
		carrier := memoryCarrier{}
		if err := codec.Inject(carrier, values); err != nil {
			b.Fatal(err)
		}
		if _, err := codec.Extract(carrier); err != nil {
			b.Fatal(err)
		}
	}
}
