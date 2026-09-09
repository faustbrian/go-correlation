package correlationotel_test

import (
	"reflect"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
	correlationotel "github.com/faustbrian/go-correlation/adapters/otel"
	legacy "github.com/faustbrian/go-correlation/telemetry"
)

func TestCanonicalOTelMatchesLegacyProjection(t *testing.T) {
	t.Parallel()

	values := correlation.Values{CorrelationID: "flow", RequestID: "request"}
	policy := correlation.DisclosurePolicy{Mode: correlation.ExposeDisclosure}
	got, gotErr := correlationotel.Attributes(values, policy)
	want, wantErr := legacy.Attributes(values, policy)
	if !reflect.DeepEqual(got, want) || gotErr != wantErr {
		t.Fatalf("Attributes() = (%#v, %v), want (%#v, %v)", got, gotErr, want, wantErr)
	}
	if correlationotel.ErrInvalidSpanContext != legacy.ErrInvalidSpanContext {
		t.Fatal("error identity differs")
	}
	if !reflect.DeepEqual(correlationotel.MetricAttributes(values), legacy.MetricAttributes(values)) {
		t.Fatal("metric attributes differ")
	}
}
