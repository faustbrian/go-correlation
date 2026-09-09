// Package telemetry links correlation metadata to OpenTelemetry without
// treating correlation IDs as trace or span IDs.
//
// Deprecated: use github.com/faustbrian/go-correlation/adapters/otel. This
// package remains supported through the documented compatibility interval.
package telemetry

import (
	correlation "github.com/faustbrian/go-correlation"
	correlationotel "github.com/faustbrian/go-correlation/adapters/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ErrInvalidSpanContext reports a missing telemetry-owned trace identity.
var ErrInvalidSpanContext = correlationotel.ErrInvalidSpanContext

// Attributes returns bounded span/link attributes under an explicit
// disclosure policy.
func Attributes(values correlation.Values, policy correlation.DisclosurePolicy) ([]attribute.KeyValue, error) {
	return correlationotel.Attributes(values, policy)
}

// Link attaches correlation attributes to a telemetry-owned SpanContext. The
// supplied trace and span IDs are neither derived from nor replaced by IDs.
func Link(spanContext trace.SpanContext, values correlation.Values, policy correlation.DisclosurePolicy) (trace.Link, error) {
	return correlationotel.Link(spanContext, values, policy)
}

// MetricAttributes exposes fixed-cardinality presence signals only. Raw,
// hashed, and deterministic business-derived identifiers are never included.
func MetricAttributes(values correlation.Values) []attribute.KeyValue {
	return correlationotel.MetricAttributes(values)
}
