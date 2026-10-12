// Package correlationslog provides at most three diagnostic slog attributes
// through Attrs, without installing a logger or handler. Disclosure defaults to
// redaction; raw/hash values stay caller-reviewed and must not become metrics.
// Calls own returned slices but borrow disclosure keys for synchronous use.
package correlationslog

import (
	"log/slog"

	correlation "github.com/faustbrian/go-correlation"
)

// Attrs returns only non-empty identifier attributes under an explicit
// disclosure policy, ordered correlation, request, causation. It returns a new
// slice or nil on error, without partial attributes. All-empty Values returns an
// empty slice without validating policy; exposed text is not validated here.
func Attrs(values correlation.Values, policy correlation.DisclosurePolicy) ([]slog.Attr, error) {
	inputs := []struct {
		key   string
		value string
	}{
		{"correlation.id", values.CorrelationID.String()},
		{"request.id", values.RequestID.String()},
		{"causation.id", values.CausationID.String()},
	}
	attributes := make([]slog.Attr, 0, len(inputs))
	for _, input := range inputs {
		if input.value == "" {
			continue
		}
		value, err := correlation.Disclose(input.key, input.value, policy)
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, slog.String(input.key, value))
	}
	return attributes, nil
}
