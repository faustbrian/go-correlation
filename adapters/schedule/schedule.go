// Package correlationschedule is the target-oriented path for scheduled-work correlation.
// Its public types and behavior retain the released schedule package identities.
package correlationschedule

import (
	correlation "github.com/faustbrian/go-correlation"
	legacy "github.com/faustbrian/go-correlation/schedule"
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
