package correlationqueue_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises the deprecated package.

import (
	"reflect"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
	queue "github.com/faustbrian/go-correlation/adapters/queue"
	legacy "github.com/faustbrian/go-correlation/queue" //nolint:staticcheck // The test proves released type identity.
)

func TestCanonicalQueuePreservesLegacyIdentity(t *testing.T) {
	t.Parallel()

	if reflect.TypeOf((*queue.Adapter)(nil)) != reflect.TypeOf((*legacy.Adapter)(nil)) {
		t.Fatal("adapter type identity differs")
	}
	if reflect.TypeOf(queue.Options{}) != reflect.TypeOf(legacy.Options{}) {
		t.Fatal("options type identity differs")
	}
	if queue.ErrInvalidOptions != legacy.ErrInvalidOptions { //nolint:errorlint // Exact sentinel identity is the contract.
		t.Fatal("error identity differs")
	}
	if got := reflect.TypeOf(queue.Options{}).PkgPath(); got != "github.com/faustbrian/go-correlation/queue" {
		t.Fatalf("options package path = %q", got)
	}

	factory, _ := correlation.NewFactory(correlation.FactoryOptions{})
	adapter, err := queue.New(factory, queue.Options{})
	if err != nil {
		t.Fatal(err)
	}
	metadata := map[string]string{}
	parent := correlation.Values{CorrelationID: "flow", RequestID: "parent"}
	if _, err := adapter.Send(metadata, parent); err != nil || metadata[correlation.DefaultCorrelationField] != "flow" {
		t.Fatalf("Send() metadata = %v, error = %v", metadata, err)
	}
}
