// Package schedule provides explicit correlation lifecycle helpers for
// scheduled work. Independent invocations start independent workflows unless
// application-owned metadata is deliberately propagated. New composes Create,
// Enqueue and Receive over an explicit Factory/queue codec. It owns no scheduler,
// job execution, acknowledgment, retry or background work; metadata and generator
// synchronization remain with their callers.
//
// Deprecated: use github.com/faustbrian/go-correlation/adapters/schedule. This
// package remains supported through the documented compatibility interval.
package schedule

import (
	"errors"

	correlation "github.com/faustbrian/go-correlation"
	queuecorrelation "github.com/faustbrian/go-correlation/adapters/queue"
)

// ErrInvalidOptions reports missing scheduled-work dependencies.
var ErrInvalidOptions = errors.New("schedule correlation: invalid options")

// Options configure optional scheduled metadata propagation.
type Options struct {
	// Queue configures metadata names/validation used by Enqueue and Receive;
	// Create uses only the supplied Factory.
	Queue queuecorrelation.Options
}

// Adapter starts independent identity roots and propagates scheduled metadata,
// without starting jobs. Use New; it inherits Factory/queue ownership rules.
type Adapter struct {
	factory *correlation.Factory
	queue   *queuecorrelation.Adapter
}

// New constructs a scheduled-work adapter retaining factory. Nil factory is
// ErrInvalidOptions; queue codec errors are preserved. Nothing is scheduled.
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

// Enqueue creates child scheduled-work metadata in the supplied nonnil map,
// refusing populated selected fields. It performs no enqueue/persistence; errors
// return zero Values and may consume generation state before injection fails.
func (adapter *Adapter) Enqueue(metadata map[string]string, parent correlation.Values) (correlation.Values, error) {
	if adapter == nil || adapter.queue == nil {
		return correlation.Values{}, ErrInvalidOptions
	}
	return adapter.queue.Send(metadata, parent)
}

// Receive receives scheduler metadata as a fresh attempt. Trust must be decided
// outside the adapter; nil metadata means absence. It delegates queue receipt
// validation/causation semantics and neither mutates the map nor executes a job.
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
