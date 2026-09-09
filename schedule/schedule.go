// Package schedule provides explicit correlation lifecycle helpers for
// scheduled work. Independent invocations start independent workflows unless
// application-owned metadata is deliberately propagated.
//
// Deprecated: use github.com/faustbrian/go-correlation/adapters/schedule. This
// package remains supported through the documented compatibility interval.
package schedule

import (
	"errors"

	correlation "github.com/faustbrian/go-correlation"
	queuecorrelation "github.com/faustbrian/go-correlation/queue"
)

// ErrInvalidOptions reports missing scheduled-work dependencies.
var ErrInvalidOptions = errors.New("schedule correlation: invalid options")

// Options configure optional scheduled metadata propagation.
type Options struct {
	Queue queuecorrelation.Options
}

// Adapter starts independent runs and explicitly propagates scheduled jobs.
type Adapter struct {
	factory *correlation.Factory
	queue   *queuecorrelation.Adapter
}

// New constructs a scheduled-work adapter.
func New(factory *correlation.Factory, options Options) (*Adapter, error) {
	if factory == nil {
		return nil, ErrInvalidOptions
	}
	queueAdapter, err := queuecorrelation.New(factory, options.Queue)
	if err != nil {
		return nil, err
	}
	return &Adapter{factory: factory, queue: queueAdapter}, nil
}

// Create begins an independent scheduler invocation with no implicit derived
// or stable workflow identity.
func (adapter *Adapter) Create() (correlation.Values, error) {
	if adapter == nil || adapter.factory == nil {
		return correlation.Values{}, ErrInvalidOptions
	}
	return adapter.factory.Create()
}

// Start begins an independent scheduler invocation with no implicit derived
// or stable workflow identity.
//
// Deprecated: use Create.
func (adapter *Adapter) Start() (correlation.Values, error) {
	return adapter.Create()
}

// Enqueue creates a child scheduled-work message.
func (adapter *Adapter) Enqueue(metadata map[string]string, parent correlation.Values) (correlation.Values, error) {
	if adapter == nil || adapter.queue == nil {
		return correlation.Values{}, ErrInvalidOptions
	}
	return adapter.queue.Send(metadata, parent)
}

// Receive receives explicitly trusted scheduler metadata as a fresh attempt.
func (adapter *Adapter) Receive(metadata map[string]string, trusted bool) (correlation.Values, error) {
	if adapter == nil || adapter.queue == nil {
		return correlation.Values{}, ErrInvalidOptions
	}
	return adapter.queue.Receive(metadata, trusted)
}

// Run receives explicitly trusted scheduler metadata as a fresh attempt.
//
// Deprecated: use Receive.
func (adapter *Adapter) Run(metadata map[string]string, trusted bool) (correlation.Values, error) {
	return adapter.Receive(metadata, trusted)
}
