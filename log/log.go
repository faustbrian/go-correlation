// Package log provides bounded slog attributes compatible with log's
// standard log/slog composition API. Attrs forwards to adapters/slog and owns no
// logger, handler, background resource or shutdown. Returned slices are owned;
// disclosure keys are borrowed, and direct raw exposure does not validate IDs.
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
// disclosure policy. It preserves adapters/slog ordering, errors, empty-value
// skipping and newly allocated output; it does not emit a log record.
func Attrs(values correlation.Values, policy correlation.DisclosurePolicy) ([]slog.Attr, error) {
	return correlationslog.Attrs(values, policy)
}
