// Package correlationjsonrpc is the target-oriented path for JSON-RPC correlation.
// Its public types and behavior retain the released jsonrpc package identities.
package correlationjsonrpc

//lint:file-ignore SA1019 The successor imports the deprecated package to preserve released type identity.

import (
	correlation "github.com/faustbrian/go-correlation"
	legacy "github.com/faustbrian/go-correlation/jsonrpc" //nolint:staticcheck // The successor must retain released type identity.
)

const (
	// CorrelationField carries the logical workflow identifier.
	CorrelationField = legacy.CorrelationField
	// RequestField carries the current JSON-RPC hop identifier.
	RequestField = legacy.RequestField
	// CausationField carries the immediate parent identifier.
	CausationField = legacy.CausationField
)

var (
	// ErrInvalidOptions reports incomplete adapter configuration.
	ErrInvalidOptions = legacy.ErrInvalidOptions
	// ErrMalformedMetadata reports non-string or oversized JSON metadata.
	ErrMalformedMetadata = legacy.ErrMalformedMetadata
)

// Metadata preserves repeated JSON members supplied by a strict envelope
// decoder so conflicting duplicates cannot be silently collapsed.
type Metadata = legacy.Metadata

// Options configure JSON-RPC metadata validation and names.
type Options = legacy.Options

// Adapter explicitly sends and receives JSON-RPC metadata.
type Adapter = legacy.Adapter

// New constructs a JSON-RPC metadata adapter.
func New(factory *correlation.Factory, options Options) (*Adapter, error) {
	return legacy.New(factory, options)
}
