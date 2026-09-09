// Package log provides bounded slog attributes compatible with log's
// standard log/slog composition API.
//
// Deprecated: use github.com/faustbrian/go-correlation/adapters/slog. This
// package remains supported through the documented compatibility interval.
package log

import (
	"log/slog"

	correlation "github.com/faustbrian/go-correlation"
	correlationslog "github.com/faustbrian/go-correlation/adapters/slog"
)

// Attrs returns only non-empty identifier attributes under an explicit
// disclosure policy.
func Attrs(values correlation.Values, policy correlation.DisclosurePolicy) ([]slog.Attr, error) {
	return correlationslog.Attrs(values, policy)
}
