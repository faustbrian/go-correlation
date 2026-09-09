package correlationjsonrpc_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises the deprecated package.

import (
	"encoding/json"
	"reflect"
	"testing"

	correlation "github.com/faustbrian/go-correlation"
	jsonrpc "github.com/faustbrian/go-correlation/adapters/jsonrpc"
	legacy "github.com/faustbrian/go-correlation/jsonrpc" //nolint:staticcheck // The test proves released type identity.
)

func TestCanonicalJSONRPCPreservesLegacyIdentity(t *testing.T) {
	t.Parallel()

	for name, types := range map[string][2]reflect.Type{
		"adapter":  {reflect.TypeOf((*jsonrpc.Adapter)(nil)), reflect.TypeOf((*legacy.Adapter)(nil))},
		"metadata": {reflect.TypeOf(jsonrpc.Metadata{}), reflect.TypeOf(legacy.Metadata{})},
		"options":  {reflect.TypeOf(jsonrpc.Options{}), reflect.TypeOf(legacy.Options{})},
	} {
		if types[0] != types[1] {
			t.Fatalf("%s type identity differs", name)
		}
	}
	if jsonrpc.ErrInvalidOptions != legacy.ErrInvalidOptions || //nolint:errorlint // Exact sentinel identity is the contract.
		jsonrpc.ErrMalformedMetadata != legacy.ErrMalformedMetadata { //nolint:errorlint // Exact sentinel identity is the contract.
		t.Fatal("error identity differs")
	}
	if got := reflect.TypeOf(jsonrpc.Metadata{}).PkgPath(); got != "github.com/faustbrian/go-correlation/jsonrpc" {
		t.Fatalf("metadata package path = %q", got)
	}

	factory, _ := correlation.NewFactory(correlation.FactoryOptions{})
	adapter, err := jsonrpc.New(factory, jsonrpc.Options{})
	if err != nil {
		t.Fatal(err)
	}
	metadata := jsonrpc.Metadata{"application": {json.RawMessage(`"kept"`)}}
	parent := correlation.Values{CorrelationID: "flow", RequestID: "parent"}
	if _, err := adapter.Send(metadata, parent); err != nil || string(metadata["application"][0]) != `"kept"` {
		t.Fatalf("Send() metadata = %v, error = %v", metadata, err)
	}
}
