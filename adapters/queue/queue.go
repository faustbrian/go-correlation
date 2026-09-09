// Package correlationqueue is the target-oriented path for queue correlation propagation.
// Its public types and behavior retain the released queue package identities.
package correlationqueue

import (
	correlation "github.com/faustbrian/go-correlation"
	legacy "github.com/faustbrian/go-correlation/queue"
)

// ErrInvalidOptions reports invalid queue propagation configuration.
var ErrInvalidOptions = legacy.ErrInvalidOptions

// Options configure queue metadata validation and field names.
type Options = legacy.Options

// Adapter explicitly sends and receives queue metadata.
type Adapter = legacy.Adapter

// New constructs a queue adapter.
func New(factory *correlation.Factory, options Options) (*Adapter, error) {
	return legacy.New(factory, options)
}
