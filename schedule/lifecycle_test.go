package schedule_test

import (
	"errors"
	"maps"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
	schedulecorrelation "github.com/faustbrian/go-correlation/schedule"
)

type lifecycleGenerator struct {
	values []string
	calls  int
	failAt int
	err    error
}

func (generator *lifecycleGenerator) New() (string, error) {
	if generator.calls == generator.failAt {
		generator.calls++
		return "", generator.err
	}
	value := generator.values[generator.calls]
	generator.calls++
	return value, nil
}

func newLifecycleAdapter(t *testing.T, values ...string) (*schedulecorrelation.Adapter, *lifecycleGenerator) {
	t.Helper()
	generator := &lifecycleGenerator{values: values, failAt: -1}
	factory, err := correlation.NewFactory(correlation.FactoryOptions{Generator: generator})
	if err != nil {
		t.Fatal(err)
	}
	adapter, err := schedulecorrelation.New(factory, schedulecorrelation.Options{})
	if err != nil {
		t.Fatal(err)
	}
	return adapter, generator
}

func newFailingLifecycleAdapter(t *testing.T, values []string, failAt int, failure error) (*schedulecorrelation.Adapter, *lifecycleGenerator) {
	t.Helper()
	generator := &lifecycleGenerator{values: values, failAt: failAt, err: failure}
	factory, err := correlation.NewFactory(correlation.FactoryOptions{Generator: generator})
	if err != nil {
		t.Fatal(err)
	}
	adapter, err := schedulecorrelation.New(factory, schedulecorrelation.Options{})
	if err != nil {
		t.Fatal(err)
	}
	return adapter, generator
}

func TestCreateAndReceiveRejectNilAdapters(t *testing.T) {
	created, createErr := (*schedulecorrelation.Adapter)(nil).Create()
	started, startErr := (*schedulecorrelation.Adapter)(nil).Start()
	if created != started || !errors.Is(createErr, schedulecorrelation.ErrInvalidOptions) || !errors.Is(startErr, schedulecorrelation.ErrInvalidOptions) || createErr.Error() != startErr.Error() {
		t.Fatalf("Create() = %#v, %v; Start() = %#v, %v", created, createErr, started, startErr)
	}

	received, receiveErr := (*schedulecorrelation.Adapter)(nil).Receive(nil, false)
	run, runErr := (*schedulecorrelation.Adapter)(nil).Run(nil, false)
	if received != run || !errors.Is(receiveErr, schedulecorrelation.ErrInvalidOptions) || !errors.Is(runErr, schedulecorrelation.ErrInvalidOptions) || receiveErr.Error() != runErr.Error() {
		t.Fatalf("Receive() = %#v, %v; Run() = %#v, %v", received, receiveErr, run, runErr)
	}
}

func TestStartDelegatesToCreate(t *testing.T) {
	createAdapter, createGenerator := newLifecycleAdapter(t, "correlation", "request")
	startAdapter, startGenerator := newLifecycleAdapter(t, "correlation", "request")

	created, createErr := createAdapter.Create()
	started, startErr := startAdapter.Start()
	want := correlation.Values{CorrelationID: "correlation", RequestID: "request"}
	if created != want || started != want || createErr != nil || startErr != nil || createGenerator.calls != 2 || startGenerator.calls != 2 {
		t.Fatalf("Create() = %#v, %v, %d; Start() = %#v, %v, %d", created, createErr, createGenerator.calls, started, startErr, startGenerator.calls)
	}
}

func TestStartDelegatesCreateFailures(t *testing.T) {
	failure := errors.New("generation")
	for _, test := range []struct {
		name   string
		failAt int
	}{
		{name: "correlation", failAt: 0},
		{name: "request", failAt: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			createAdapter, createGenerator := newFailingLifecycleAdapter(t, []string{"correlation"}, test.failAt, failure)
			startAdapter, startGenerator := newFailingLifecycleAdapter(t, []string{"correlation"}, test.failAt, failure)

			created, createErr := createAdapter.Create()
			started, startErr := startAdapter.Start()
			wantCalls := test.failAt + 1
			if created != (correlation.Values{}) || started != (correlation.Values{}) ||
				!errors.Is(createErr, correlation.ErrGeneration) || !errors.Is(startErr, correlation.ErrGeneration) ||
				!errors.Is(createErr, failure) || !errors.Is(startErr, failure) ||
				createErr.Error() != startErr.Error() || createGenerator.calls != wantCalls || startGenerator.calls != wantCalls {
				t.Fatalf("Create() = %#v, %v, %d; Start() = %#v, %v, %d", created, createErr, createGenerator.calls, started, startErr, startGenerator.calls)
			}
		})
	}
}

func TestRunDelegatesToReceiveWithoutMutatingMetadata(t *testing.T) {
	for _, test := range []struct {
		name     string
		metadata map[string]string
		trusted  bool
		values   []string
		want     correlation.Values
		wantErrs []error
		calls    int
		failAt   int
		failure  bool
	}{
		{name: "absent", metadata: nil, values: []string{"correlation", "request"}, want: correlation.Values{CorrelationID: "correlation", RequestID: "request"}, calls: 2},
		{name: "untrusted", metadata: map[string]string{correlation.DefaultCorrelationField: "inbound", correlation.DefaultRequestField: "parent"}, values: []string{"correlation", "request"}, want: correlation.Values{CorrelationID: "correlation", RequestID: "request"}, calls: 2},
		{name: "untrusted malformed", metadata: map[string]string{correlation.DefaultCorrelationField: "bad value"}, wantErrs: []error{correlation.ErrInvalidCarrier}},
		{name: "trusted", metadata: map[string]string{correlation.DefaultCorrelationField: "correlation", correlation.DefaultRequestField: "parent"}, trusted: true, values: []string{"request"}, want: correlation.Values{CorrelationID: "correlation", RequestID: "request", CausationID: "parent"}, calls: 1},
		{name: "trusted malformed", metadata: map[string]string{correlation.DefaultCorrelationField: "bad value"}, trusted: true, wantErrs: []error{correlation.ErrInvalidCarrier}},
		{name: "untrusted correlation generation failure", metadata: map[string]string{correlation.DefaultCorrelationField: "inbound", correlation.DefaultRequestField: "parent"}, wantErrs: []error{correlation.ErrGeneration, errLifecycleGeneration}, calls: 1, failAt: 0, failure: true},
		{name: "untrusted request generation failure", metadata: map[string]string{correlation.DefaultCorrelationField: "inbound", correlation.DefaultRequestField: "parent"}, values: []string{"correlation"}, wantErrs: []error{correlation.ErrGeneration, errLifecycleGeneration}, calls: 2, failAt: 1, failure: true},
		{name: "trusted request generation failure", metadata: map[string]string{correlation.DefaultCorrelationField: "correlation", correlation.DefaultRequestField: "parent"}, trusted: true, wantErrs: []error{correlation.ErrGeneration, errLifecycleGeneration}, calls: 1, failAt: 0, failure: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var receiveAdapter, runAdapter *schedulecorrelation.Adapter
			var receiveGenerator, runGenerator *lifecycleGenerator
			if test.failure {
				receiveAdapter, receiveGenerator = newFailingLifecycleAdapter(t, test.values, test.failAt, errLifecycleGeneration)
				runAdapter, runGenerator = newFailingLifecycleAdapter(t, test.values, test.failAt, errLifecycleGeneration)
			} else {
				receiveAdapter, receiveGenerator = newLifecycleAdapter(t, test.values...)
				runAdapter, runGenerator = newLifecycleAdapter(t, test.values...)
			}
			receiveMetadata := maps.Clone(test.metadata)
			runMetadata := maps.Clone(test.metadata)

			received, receiveErr := receiveAdapter.Receive(receiveMetadata, test.trusted)
			run, runErr := runAdapter.Run(runMetadata, test.trusted)
			if received != test.want || received != run || !sameError(receiveErr, runErr) ||
				!matchesErrors(receiveErr, test.wantErrs) || !matchesErrors(runErr, test.wantErrs) ||
				receiveGenerator.calls != test.calls || runGenerator.calls != test.calls {
				t.Fatalf("Receive() = %#v, %v, %d; Run() = %#v, %v, %d", received, receiveErr, receiveGenerator.calls, run, runErr, runGenerator.calls)
			}
			if !maps.Equal(receiveMetadata, test.metadata) || !maps.Equal(runMetadata, test.metadata) {
				t.Fatalf("metadata mutated: Receive = %#v, Run = %#v", receiveMetadata, runMetadata)
			}
		})
	}
}

var errLifecycleGeneration = errors.New("lifecycle generation")

func matchesErrors(got error, wants []error) bool {
	if len(wants) == 0 {
		return got == nil
	}
	for _, want := range wants {
		if !errors.Is(got, want) {
			return false
		}
	}
	return true
}

func sameError(first, second error) bool {
	if first == nil {
		return second == nil
	}
	if second == nil {
		return false
	}
	return first.Error() == second.Error()
}
