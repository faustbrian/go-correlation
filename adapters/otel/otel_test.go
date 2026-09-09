package correlationotel_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises the deprecated package.

import (
	"reflect"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
	correlationotel "github.com/faustbrian/go-correlation/adapters/otel"
	legacy "github.com/faustbrian/go-correlation/telemetry" //nolint:staticcheck // The test proves released sentinel compatibility.
)

func TestCanonicalOTelMatchesLegacyProjection(t *testing.T) {
	t.Parallel()

	values := correlation.Values{CorrelationID: "flow", RequestID: "request"}
	policy := correlation.DisclosurePolicy{Mode: correlation.ExposeDisclosure}
	got, gotErr := correlationotel.Attributes(values, policy)
	want, wantErr := legacy.Attributes(values, policy)
	if !reflect.DeepEqual(got, want) || gotErr != wantErr { //nolint:errorlint // Exact error identity proves delegation parity.
		t.Fatalf("Attributes() = (%#v, %v), want (%#v, %v)", got, gotErr, want, wantErr)
	}
	if correlationotel.ErrInvalidSpanContext != legacy.ErrInvalidSpanContext { //nolint:errorlint // Exact sentinel identity is the contract.
		t.Fatal("error identity differs")
	}
	if !reflect.DeepEqual(correlationotel.MetricAttributes(values), legacy.MetricAttributes(values)) {
		t.Fatal("metric attributes differ")
	}
}
