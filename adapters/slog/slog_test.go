package correlationslog_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises the deprecated package.

import (
	"reflect"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
	correlationslog "github.com/faustbrian/go-correlation/adapters/slog"
	legacy "github.com/faustbrian/go-correlation/log" //nolint:staticcheck // The test proves released projection compatibility.
)

func TestCanonicalSlogMatchesLegacyProjection(t *testing.T) {
	t.Parallel()

	values := correlation.Values{CorrelationID: "flow", RequestID: "request"}
	policy := correlation.DisclosurePolicy{Mode: correlation.ExposeDisclosure}
	got, gotErr := correlationslog.Attrs(values, policy)
	want, wantErr := legacy.Attrs(values, policy)
	if !reflect.DeepEqual(got, want) || gotErr != wantErr { //nolint:errorlint // Exact error identity proves delegation parity.
		t.Fatalf("Attrs() = (%#v, %v), want (%#v, %v)", got, gotErr, want, wantErr)
	}
}
