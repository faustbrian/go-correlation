// Package correlationschedule is the target-oriented path for scheduled-work correlation.
// Its public types and behavior retain the released schedule package identities.
package correlationschedule

//lint:file-ignore SA1019 The successor imports the deprecated package to preserve released type identity.

import (
	correlation "github.com/faustbrian/go-correlation"
	legacy "github.com/faustbrian/go-correlation/schedule" //nolint:staticcheck // The successor must retain released type identity.
)

// ErrInvalidOptions reports missing scheduled-work dependencies.
var ErrInvalidOptions = legacy.ErrInvalidOptions

// Options configure optional scheduled metadata propagation.
type Options = legacy.Options

// Adapter starts independent runs and explicitly propagates scheduled jobs.
type Adapter = legacy.Adapter

// New constructs a scheduled-work adapter.
func New(factory *correlation.Factory, options Options) (*Adapter, error) {
	return legacy.New(factory, options)
}
