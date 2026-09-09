package httpcorrelation_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises the deprecated package.

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
	httpcorrelation "github.com/faustbrian/go-correlation/adapters/http"
	legacy "github.com/faustbrian/go-correlation/http" //nolint:staticcheck // The test proves released type identity.
)

func TestCanonicalHTTPPreservesLegacyIdentity(t *testing.T) {
	t.Parallel()

	if reflect.TypeOf((*httpcorrelation.Middleware)(nil)) != reflect.TypeOf((*legacy.Middleware)(nil)) {
		t.Fatal("middleware type identity differs")
	}
	if reflect.TypeOf(httpcorrelation.Options{}) != reflect.TypeOf(legacy.Options{}) {
		t.Fatal("options type identity differs")
	}
	if httpcorrelation.ErrInvalidOptions != legacy.ErrInvalidOptions { //nolint:errorlint // Exact sentinel identity is the contract.
		t.Fatal("error identity differs")
	}
	if got := reflect.TypeOf(httpcorrelation.Options{}).PkgPath(); got != "github.com/faustbrian/go-correlation/http" {
		t.Fatalf("options package path = %q", got)
	}

	factory, _ := correlation.NewFactory(correlation.FactoryOptions{})
	middleware, err := httpcorrelation.New(factory, httpcorrelation.Options{})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	parent := correlation.Values{CorrelationID: "flow", RequestID: "parent"}
	if _, err := middleware.Inject(request, parent); err != nil || request.Header.Get(httpcorrelation.CorrelationHeader) != "flow" {
		t.Fatalf("Inject() headers = %v, error = %v", request.Header, err)
	}
}
