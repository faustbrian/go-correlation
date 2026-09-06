package correlation_test

import (
	"errors"
	"strings"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
)

type generatorResult struct {
	value string
	err   error
}

type recordingGenerator struct {
	results []generatorResult
	calls   int
}

type mapGenerator map[string]string

func (generator mapGenerator) New() (string, error) { return generator["value"], nil }

type sliceGenerator []string

func (generator sliceGenerator) New() (string, error) { return generator[0], nil }

type channelGenerator chan string

func (generator channelGenerator) New() (string, error) { return <-generator, nil }

func (generator *recordingGenerator) New() (string, error) {
	result := generator.results[generator.calls]
	generator.calls++
	return result.value, result.err
}

func TestNewFactoryRejectsTypedNilGenerators(t *testing.T) {
	var pointer *recordingGenerator
	var function correlation.GeneratorFunc
	var mapped mapGenerator
	var sliced sliceGenerator
	var channel channelGenerator

	for name, generator := range map[string]correlation.Generator{
		"channel":  channel,
		"function": function,
		"map":      mapped,
		"pointer":  pointer,
		"slice":    sliced,
	} {
		t.Run(name, func(t *testing.T) {
			factory, err := correlation.NewFactory(correlation.FactoryOptions{Generator: generator})
			if factory != nil || !errors.Is(err, correlation.ErrInvalidFactory) {
				t.Fatalf("NewFactory() = %#v, %v", factory, err)
			}
		})
	}
}

func TestNewFactoryPreservesPolicyValidationPrecedence(t *testing.T) {
	var generator *recordingGenerator
	factory, err := correlation.NewFactory(correlation.FactoryOptions{
		Policy:    correlation.Policy{MaxLength: -1},
		Generator: generator,
	})
	if factory != nil || !errors.Is(err, correlation.ErrInvalidFactory) || !strings.Contains(err.Error(), "maximum length") {
		t.Fatalf("NewFactory() = %#v, %v", factory, err)
	}
}

func TestFactoryCreateGeneratesCorrelationThenRequest(t *testing.T) {
	generator := &recordingGenerator{results: []generatorResult{{value: "correlation"}, {value: "request"}}}
	factory, err := correlation.NewFactory(correlation.FactoryOptions{Generator: generator})
	if err != nil {
		t.Fatal(err)
	}

	values, err := factory.Create()
	if err != nil {
		t.Fatal(err)
	}
	if values != (correlation.Values{CorrelationID: "correlation", RequestID: "request"}) || generator.calls != 2 {
		t.Fatalf("Create() = %#v, calls = %d", values, generator.calls)
	}
}

func TestFactoryCreateReturnsNoPartialValuesOnGenerationFailure(t *testing.T) {
	firstFailure := errors.New("first generation")
	secondFailure := errors.New("second generation")

	for name, test := range map[string]struct {
		results []generatorResult
		want    error
		calls   int
	}{
		"correlation": {results: []generatorResult{{err: firstFailure}}, want: firstFailure, calls: 1},
		"request":     {results: []generatorResult{{value: "correlation"}, {err: secondFailure}}, want: secondFailure, calls: 2},
	} {
		t.Run(name, func(t *testing.T) {
			generator := &recordingGenerator{results: test.results}
			factory, err := correlation.NewFactory(correlation.FactoryOptions{Generator: generator})
			if err != nil {
				t.Fatal(err)
			}

			values, err := factory.Create()
			if values != (correlation.Values{}) || !errors.Is(err, correlation.ErrGeneration) || !errors.Is(err, test.want) || generator.calls != test.calls {
				t.Fatalf("Create() = %#v, %v, calls = %d", values, err, generator.calls)
			}
		})
	}
}

func TestFactoryStartDelegatesToCreate(t *testing.T) {
	createGenerator := &recordingGenerator{results: []generatorResult{{value: "correlation"}, {value: "request"}}}
	startGenerator := &recordingGenerator{results: []generatorResult{{value: "correlation"}, {value: "request"}}}
	createFactory, _ := correlation.NewFactory(correlation.FactoryOptions{Generator: createGenerator})
	startFactory, _ := correlation.NewFactory(correlation.FactoryOptions{Generator: startGenerator})

	created, createErr := createFactory.Create()
	started, startErr := startFactory.Start()
	if created != started || createErr != nil || startErr != nil || createGenerator.calls != startGenerator.calls {
		t.Fatalf("Create() = %#v, %v, %d; Start() = %#v, %v, %d", created, createErr, createGenerator.calls, started, startErr, startGenerator.calls)
	}

	created, createErr = (*correlation.Factory)(nil).Create()
	started, startErr = (*correlation.Factory)(nil).Start()
	if created != started || !errors.Is(createErr, correlation.ErrGeneration) || !errors.Is(startErr, correlation.ErrGeneration) || createErr.Error() != startErr.Error() {
		t.Fatalf("nil Create() = %#v, %v; nil Start() = %#v, %v", created, createErr, started, startErr)
	}
}

func TestFactoryStartDelegatesGenerationFailures(t *testing.T) {
	firstFailure := errors.New("first generation")
	secondFailure := errors.New("second generation")

	for name, test := range map[string]struct {
		results []generatorResult
		want    error
		calls   int
	}{
		"correlation": {results: []generatorResult{{err: firstFailure}}, want: firstFailure, calls: 1},
		"request":     {results: []generatorResult{{value: "correlation"}, {err: secondFailure}}, want: secondFailure, calls: 2},
	} {
		t.Run(name, func(t *testing.T) {
			createGenerator := &recordingGenerator{results: test.results}
			startGenerator := &recordingGenerator{results: test.results}
			createFactory, _ := correlation.NewFactory(correlation.FactoryOptions{Generator: createGenerator})
			startFactory, _ := correlation.NewFactory(correlation.FactoryOptions{Generator: startGenerator})

			created, createErr := createFactory.Create()
			started, startErr := startFactory.Start()
			if created != started ||
				!errors.Is(createErr, correlation.ErrGeneration) || !errors.Is(startErr, correlation.ErrGeneration) ||
				!errors.Is(createErr, test.want) || !errors.Is(startErr, test.want) ||
				createErr.Error() != startErr.Error() || createGenerator.calls != test.calls || startGenerator.calls != test.calls {
				t.Fatalf("Create() = %#v, %v, %d; Start() = %#v, %v, %d", created, createErr, createGenerator.calls, started, startErr, startGenerator.calls)
			}
		})
	}
}
