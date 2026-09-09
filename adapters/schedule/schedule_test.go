package correlationschedule_test

import (
	"reflect"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
	schedule "github.com/faustbrian/go-correlation/adapters/schedule"
	legacy "github.com/faustbrian/go-correlation/schedule"
)

func TestCanonicalSchedulePreservesLegacyIdentity(t *testing.T) {
	t.Parallel()

	if reflect.TypeOf((*schedule.Adapter)(nil)) != reflect.TypeOf((*legacy.Adapter)(nil)) {
		t.Fatal("adapter type identity differs")
	}
	if reflect.TypeOf(schedule.Options{}) != reflect.TypeOf(legacy.Options{}) {
		t.Fatal("options type identity differs")
	}
	if schedule.ErrInvalidOptions != legacy.ErrInvalidOptions {
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
