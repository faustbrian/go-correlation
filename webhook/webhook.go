// Package webhook provides an explicit webhook-named HTTP correlation adapter.
// New composes Send and Wrap over HTTP policy without webhook signing, delivery,
// body access, retry or replay protection. Peer trust, resource lifetime and
// request/header synchronization remain application-owned.
package webhook

import (
	"net/http"

	correlation "github.com/faustbrian/go-correlation"
	httpcorrelation "github.com/faustbrian/go-correlation/adapters/http"
)

// Options configure webhook boundary trust and invalid input handling.
type Options = httpcorrelation.Options

// Adapter keeps webhook propagation explicit at send and receive call sites.
// Construct with New; nil Adapter receivers panic on Send/Wrap. Sharing inherits
// HTTP Factory/Trust callback safety and does not synchronize header maps.
type Adapter struct{ middleware *httpcorrelation.Middleware }

// New constructs an explicit webhook correlation adapter.
func New(factory *correlation.Factory, options Options) (*Adapter, error) {
	middleware, err := httpcorrelation.New(factory, options)
	if err != nil {
		return nil, err
	}
	return &Adapter{middleware: middleware}, nil
}

// Send creates and injects an outbound webhook request hop without sending it.
// It delegates HTTP Inject's header ownership, overwrite and failure semantics;
// it does not create signatures or authenticate a receiver.
func (adapter *Adapter) Send(request *http.Request, parent correlation.Values) (correlation.Values, error) {
	return adapter.middleware.Inject(request, parent)
}

// Wrap receives a webhook through the configured HTTP trust/invalid-input policy.
// It delegates HTTP Wrap's canonical header writes and derived context, not
// webhook authentication or body processing; callbacks/handlers are borrowed.
func (adapter *Adapter) Wrap(next http.Handler) http.Handler {
	return adapter.middleware.Wrap(next)
}
