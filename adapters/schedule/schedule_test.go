package correlationschedule_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises the deprecated package.

import (
	"reflect"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
	schedule "github.com/faustbrian/go-correlation/adapters/schedule"
	legacy "github.com/faustbrian/go-correlation/schedule" //nolint:staticcheck // The test proves released type identity.
)

func TestCanonicalSchedulePreservesLegacyIdentity(t *testing.T) {
	t.Parallel()

	if reflect.TypeOf((*schedule.Adapter)(nil)) != reflect.TypeOf((*legacy.Adapter)(nil)) {
		t.Fatal("adapter type identity differs")
	}
	if reflect.TypeOf(schedule.Options{}) != reflect.TypeOf(legacy.Options{}) {
		t.Fatal("options type identity differs")
	}
	if schedule.ErrInvalidOptions != legacy.ErrInvalidOptions { //nolint:errorlint // Exact sentinel identity is the contract.
		t.Fatal("error identity differs")
	}
	if got := reflect.TypeOf(schedule.Options{}).PkgPath(); got != "github.com/faustbrian/go-correlation/schedule" {
		t.Fatalf("options package path = %q", got)
	}

	factory, _ := correlation.NewFactory(correlation.FactoryOptions{})
	adapter, err := schedule.New(factory, schedule.Options{})
	if err != nil {
		t.Fatal(err)
	}
	values, err := adapter.Create()
	if err != nil || values.CorrelationID == "" || values.RequestID == "" {
		t.Fatalf("Create() = %#v, %v", values, err)
	}
}
