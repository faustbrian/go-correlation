# Operations

Monitor malformed and conflicting carrier rejections by bounded reason code,
not identifier value. A sudden rise usually indicates a proxy combining
headers, an adapter mismatch, or hostile input. Treat generator failures as
internal errors and investigate system entropy or injected generator health.

When debugging a workflow, disclose raw IDs only in access-controlled logs
with an explicit retention period. Never copy them into metric labels. During
key rotation, version deterministic domains and allow the old version only for
the bounded migration window.

## Runtime and ownership

The module starts no goroutines, opens no sockets or files, and exposes no
shutdown method. Each call completes synchronously. HTTP request cancellation
and deadlines remain owned by the caller and downstream handler; correlation
middleware only derives the request context used for that handler. Queue,
schedule, JSON-RPC, and webhook adapters do not retry or acknowledge work.

Factories, codecs, deterministic strategies, and adapters do not mutate their
configured values after construction. The built-in generator is safe for
shared use. A custom generator and HTTP `Trust` callback are borrowed, so the
caller must keep them valid for the owning object's lifetime and provide any
synchronization they require. Receive and extraction paths only read carriers;
send, inject, and HTTP wrap paths update caller-owned maps or headers, which
must not be written concurrently.

## Performance and limits

Random root creation consumes two generated identifiers; `Next` consumes one.
The built-in factory buffers 4 KiB of cryptographic entropy per factory to
reduce system reads without process-global state. Custom identifier policies
are capped at 1,024 bytes, and the default cap is 128 bytes.

Carrier parsing is linear in the small bounded field set. HTTP and JSON-RPC
adapters reject more than eight values for one field. Deterministic derivation
hashes at most 1 MiB per call and copies the configured key at construction.
Do not place correlation values in metric labels: even bounded-length IDs have
unbounded cardinality across requests.

Use the package [benchmarks](../benchmark_test.go) to compare changes under the
documented repository environment. They are engineering evidence, not an
application latency service-level objective.

For symptom-driven diagnosis, see [troubleshooting](troubleshooting.md).
