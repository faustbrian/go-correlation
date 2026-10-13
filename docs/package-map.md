# Package map

[Documentation index](README.md) · [Repository entry](../README.md)

This map follows the current source modules, not a synchronized release.
Use the repository README for published installation versions; an unpublished
source major is not installable merely because its package appears here.
Internal packages and non-releasable fixtures are intentionally excluded.

## github.com/faustbrian/go-correlation

`correlation` is the transport-neutral owner of correlation, request, and causation identifiers. It carries those semantics through HTTP, JSON-RPC, queues, scheduled work, webhooks, logs, and OpenTelemetry without creating a global propagator or redefining trace and idempotency concepts.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-correlation`](../carrier.go) | Use for immutable correlation values, validated identifiers, trust policy, disclosure policy, and transport-neutral propagation. |
| [`github.com/faustbrian/go-correlation/adapters/http`](../adapters/http/http.go) | Use for bounded inbound and outbound HTTP correlation with explicit proxy trust. |
| [`github.com/faustbrian/go-correlation/adapters/jsonrpc`](../adapters/jsonrpc/jsonrpc.go) | Use for explicit correlation metadata alongside JSON-RPC envelopes. |
| [`github.com/faustbrian/go-correlation/adapters/otel`](../adapters/otel/otel.go) | Use to attach bounded correlation attributes and links to caller-owned OpenTelemetry data. |
| [`github.com/faustbrian/go-correlation/adapters/queue`](../adapters/queue/queue.go) | Use to preserve workflow identity while rotating request identity for each queue delivery attempt. |
| [`github.com/faustbrian/go-correlation/adapters/schedule`](../adapters/schedule/schedule.go) | Use to create or deliberately propagate correlation values for scheduled runs. |
| [`github.com/faustbrian/go-correlation/adapters/slog`](../adapters/slog/slog.go) | Use for bounded redacted, keyed-hash, or explicitly raw slog attributes. |
| [`github.com/faustbrian/go-correlation/http`](../http/http.go) | Compatibility path for the canonical HTTP adapter. |
| [`github.com/faustbrian/go-correlation/http/requestidbridge`](../http/requestidbridge/requestidbridge.go) | Use to adopt a trusted go-http-middleware request identifier without sharing private context keys. |
| [`github.com/faustbrian/go-correlation/jsonrpc`](../jsonrpc/jsonrpc.go) | Compatibility path for the canonical JSON-RPC adapter. |
| [`github.com/faustbrian/go-correlation/log`](../log/log.go) | Compatibility path for the canonical slog adapter. |
| [`github.com/faustbrian/go-correlation/queue`](../queue/queue.go) | Compatibility path for the canonical queue adapter. |
| [`github.com/faustbrian/go-correlation/schedule`](../schedule/schedule.go) | Compatibility path for the canonical scheduler adapter. |
| [`github.com/faustbrian/go-correlation/telemetry`](../telemetry/telemetry.go) | Compatibility path for the canonical OpenTelemetry adapter. |
| [`github.com/faustbrian/go-correlation/webhook`](../webhook/webhook.go) | Use for explicit outbound and inbound webhook correlation over HTTP. |

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
