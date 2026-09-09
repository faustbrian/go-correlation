package correlationslog_test

import (
	"reflect"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
	correlationslog "github.com/faustbrian/go-correlation/adapters/slog"
	legacy "github.com/faustbrian/go-correlation/log"
)

func TestCanonicalSlogMatchesLegacyProjection(t *testing.T) {
	t.Parallel()

	values := correlation.Values{CorrelationID: "flow", RequestID: "request"}
	policy := correlation.DisclosurePolicy{Mode: correlation.ExposeDisclosure}
	got, gotErr := correlationslog.Attrs(values, policy)
	want, wantErr := legacy.Attrs(values, policy)
	if !reflect.DeepEqual(got, want) || gotErr != wantErr {
		t.Fatalf("Attrs() = (%#v, %v), want (%#v, %v)", got, gotErr, want, wantErr)
	}
}
