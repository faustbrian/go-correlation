// Package httpcorrelation is the target-oriented path for HTTP correlation.
// Its public types and behavior retain the released http package identities.
package httpcorrelation

import (
	correlation "github.com/faustbrian/go-correlation"
	legacy "github.com/faustbrian/go-correlation/http"
)

const (
	// CorrelationHeader carries the logical workflow identifier.
	CorrelationHeader = legacy.CorrelationHeader
	// RequestHeader carries the current HTTP hop identifier.
	RequestHeader = legacy.RequestHeader
	// CausationHeader carries the immediate parent identifier.
	CausationHeader = legacy.CausationHeader
)

// InvalidPolicy controls malformed or ambiguous inbound metadata.
type InvalidPolicy = legacy.InvalidPolicy

const (
	// ReplaceInvalid discards every inbound value and creates a fresh root.
	ReplaceInvalid = legacy.ReplaceInvalid
	// RejectInvalid returns HTTP 400 before application code runs.
	RejectInvalid = legacy.RejectInvalid
)

// ErrInvalidOptions reports invalid HTTP adapter configuration.
var ErrInvalidOptions = legacy.ErrInvalidOptions

// Options configure immutable HTTP propagation policy.
type Options = legacy.Options

// Middleware owns explicit HTTP extraction, trust, context, and injection.
type Middleware = legacy.Middleware

// New constructs an HTTP adapter without installing global middleware.
func New(factory *correlation.Factory, options Options) (*Middleware, error) {
	return legacy.New(factory, options)
}
