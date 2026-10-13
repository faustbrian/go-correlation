# Source documentation audit

[Documentation index](README.md) · [Repository entry](../README.md)

## Scope and reproduction

Tier A audit at exact remote-main base
`e0c9a2ac551aa543261d2b87c2236e12d05de298`, selected detached worktree.
Reproduce with `git ls-files '*.go'`, modules.json/packages.json and
`golib source-docs inventory --json`. All 51 tracked Go files are included:
21 production and 30 assurance sources. Two modules: public root
`github.com/faustbrian/go-correlation` and non-releasable
`github.com/faustbrian/go-correlation/integration/siblings`. The root declares
15 public packages; the nested harness has one test-only package and consumes
released public modules rather than the edited root via local replacement.
No untracked or sibling implementation is part of this source boundary.

Before source edits, objective inventory finds 37 missing exported field,
interface-method and String-method comments. All 15 public package comments
exist but need meaningful ownership/trust or compatibility boundaries.
Source investigation, distinct from missing-comment repair, covers direct
string conversions bypassing validation, default-generated UUID policy bypass,
raw disclosure without validation, ignored inbound causation, callback/cause
privacy, exact carrier preflight versus transactional guarantees, and canonical
HTTP writes versus complete header sanitation. No rationale is invented for
these implementation distinctions.

Outcome: truthful package/API contracts and adjacent invariant rationale for
root values/generation/transport/observability users. Risks are identifier
trust confusion, accidental raw disclosure, custom callbacks blocking or
panicking, mutable map/key sharing and compatibility alias identity. Preserve
all executable behavior, API signatures, dependency pins and historical records.
Select native formatting and all-file AST/token/directive equivalence; render
all packages; native full go-doc API baseline regeneration in configured order;
native docs (local links, Go docs and three examples); objective inventory,
changed-content secret checks and ordinary complete final-diff review. Runtime,
race, active fuzz, performance, sibling-service/release and consumer campaigns
are outside this prose-only pass. Go checks use one serialized owned-cache slot.

## Investigated source contracts

- IDs are named strings, not validating or secret-safe types. Empty Values fields
  mean absence; Parse rejects empty/invalid bytes without normalization, while
  direct conversions/context storage do not validate. The policy counts ASCII
  bytes, default 128, supported positive maximum 1024. Must helpers panic.
- Factory configuration is immutable but its generator is borrowed mutable code.
  Default UUIDv4 uses Identifier v2.0.1 and a 4 KiB entropy buffer; concurrency
  depends on that generator's synchronization, not a Factory-owned mutex.
  Custom callbacks must synchronize, supply suitable randomness/uniqueness and
  avoid leaking private details through wrapped errors. Nothing closes them.
  Default UUID output is trusted rather than revalidated against custom length:
  a configured limit below 36 does not shorten/refuse built-in generated IDs.
- Create consumes correlation then request generation, returning zero Values on
  error without rolling back generator state. Next validates parent correlation
  and request, ignores parent causation and generates one request. Accept ignores
  inbound causation, validates only opted-in nonempty correlation/request fields,
  and always generates a new request. Carrier extraction validates all fields
  before trust, so untrusted malformed metadata can still fail receipt.
- Carrier providers supply stable copies/immutable views and exact replace Set.
  Extract accepts absent fields or up to eight equal nonempty values; no partial
  Values on failure. Inject preflights all nonempty fields before Set, not skipped
  empty output fields; this is not atomic against concurrent carrier mutation or
  callback failure. Codec/Propagator retain no caller carrier and own no retry.
- HTTP scans header keys case-insensitively, preserving duplicates for validation.
  RejectInvalid creates fresh diagnostics then returns 400 before the handler;
  generation failure returns 500. Default ReplaceInvalid with no Trust skips
  inbound extraction entirely. Trust is a synchronous borrowed callback, never
  authentication performed by the adapter. Wrap writes canonical map keys in the
  original headers and response and derives context, not a deep request clone or
  universal removal of alternate-cased keys. Send mutates supplied headers.
- JSON-RPC only touches explicitly named metadata keys, not envelopes or protocol
  request IDs. Providers retain duplicate raw members. Receipt bounds each
  relevant encoded value at 1026 bytes and eight values per field before decoding;
  unknown metadata is ignored. Queue and schedule borrow maps; nil receive maps
  mean absent metadata, nil send maps fail. Scheduler helpers do not run jobs.
- Deterministic is HMAC-SHA-256 even with an empty key; transcript uses package
  marker, big-endian version and length-delimited domain/input. Constructor
  copies at most 1024 key bytes, validates ASCII domain <=128 and positive
  version. Output includes d<version>_ and a truncated unpadded base64url digest,
  at least 16 digest characters. Input is borrowed only for the synchronous call,
  nonempty and <=1 MiB; truncation trades collision resistance for output length.
  Empty-key low-entropy inputs are enumerable; stable outputs are linkable, not
  authentication or collision-free identities. No secure erasure API is owned.
- Disclose borrows a <=1024-byte key. Redaction is default; hashes require nonempty
  key and emit 16 HMAC digest bytes as unpadded base64url. Expose returns raw input
  without validation. Neither labels nor values are syntax/size checked here;
  use validated diagnostic IDs and package-owned labels. Empty values return
  absent output before mode validation (but after key size). Projection helpers
  skip absent Values, so invalid disclosure policy can go unnoticed on all-empty
  Values. Only metric presence booleans are fixed-cardinality.
- External trust and source are declarations by the caller, never verified.
  Context functions copy Values, not validate them; nil lookup context is absent,
  but nil WithValues parent panics through context.WithValue. Request-ID bridge
  replaces only the supplied empty request ID using explicit trusted lookup;
  other Values fields are not validated or merged from the parent context.
- Factories/strategies/codecs/adapters start no background work, network calls,
  timers or retries and own no shutdown. Custom generators/Trust/lookups/handlers
  can block or panic; no deadline or panic recovery is added. Carrier/header/maps
  and Disclose keys remain caller-synchronized. Nil/zero receiver handling is
  method-specific, not a package-wide safe-zero guarantee.

## Variants, provenance and assurance

No generated Go source, generator directive, platform suffix or build constraint
is tracked. No unfinished implementation markers were found. Existing valid
package/method deprecations retain the documented compatibility interval;
security conversion suppressions and compatibility lint suppressions remain exact.
Canonical HTTP/JSON-RPC/queue/schedule packages intentionally alias legacy types
and sentinels; canonical slog/OTel own projections and deprecated packages forward.
Changing alias direction could alter observable type/sentinel identity.

Four fuzz targets: core parse identity, carrier reinjection, HTTP status and
JSON-RPC raw byte/count boundaries. Five benchmarks use b.Loop with setup outside
measurement: oversize rejection, custom/default hop generation, default root
creation and per-iteration carrier roundtrip. Three examples demonstrate root
creation and hop semantics. Synthetic sequences/carriers are test-local seams;
committed deterministic vectors protect this package's transcript, not an
external conformance claim. The nested harness pins separate published sibling
versions and its unchanged results are not proof of this worktree's new source.
No protocol-spec register or external generated corpus is maintained here;
JSON-RPC envelope and W3C trace/baggage conformance are expressly not owned.

## Tracked Go coverage

| File | Audited boundary |
| --- | --- |
| `adapters/http/http.go` | Canonical projection or compatibility-preserving alias contracts |
| `adapters/http/http_test.go` | Canonical adapter identity/forwarding parity |
| `adapters/jsonrpc/jsonrpc.go` | Canonical projection or compatibility-preserving alias contracts |
| `adapters/jsonrpc/jsonrpc_test.go` | Canonical adapter identity/forwarding parity |
| `adapters/otel/otel.go` | Canonical projection or compatibility-preserving alias contracts |
| `adapters/otel/otel_test.go` | Canonical adapter identity/forwarding parity |
| `adapters/queue/queue.go` | Canonical projection or compatibility-preserving alias contracts |
| `adapters/queue/queue_test.go` | Canonical adapter identity/forwarding parity |
| `adapters/schedule/schedule.go` | Canonical projection or compatibility-preserving alias contracts |
| `adapters/schedule/schedule_test.go` | Canonical adapter identity/forwarding parity |
| `adapters/slog/slog.go` | Canonical projection or compatibility-preserving alias contracts |
| `adapters/slog/slog_test.go` | Canonical adapter identity/forwarding parity |
| `benchmark_test.go` | Five b.Loop workloads, excluded setup, custom versus UUID generators |
| `carrier.go` | Carrier obligations, bounded duplicates, exact field names, preflight and overwrite policy |
| `carrier_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `concurrency_test.go` | Default generator/context concurrency and immutable copied derivation key |
| `correlation.go` | Typed strings, ASCII byte validation, zero policy and context copying |
| `correlation_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `deterministic.go` | Versioned HMAC transcript, truncation, copied key and linkability |
| `deterministic_vectors_test.go` | Committed HMAC/truncation vectors, not external conformance fixtures |
| `disclosure.go` | Redaction/hash/raw modes, borrowed keys and absence versus invalid policy |
| `edge_test.go` | Validation/generation/privacy/derivation/carrier/error boundary partitions |
| `example_test.go` | Three runnable root and hop examples |
| `external.go` | Declared external kind/source/trust without implicit promotion |
| `fuzz_test.go` | Two core hostile-input boundaries and lossless/roundtrip oracles |
| `generation.go` | Borrowed generators, default UUID path, call order and partial failure |
| `generation_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `http/edge_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `http/fuzz_test.go` | Header rejection status oracle without peer trust or real network |
| `http/http.go` | Case-folded inbound headers, peer trust callback, canonical writes and response handling |
| `http/http_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `http/requestidbridge/requestidbridge.go` | Explicit trusted lookup, overwrite guard and derived context |
| `http/requestidbridge/requestidbridge_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `integration/siblings/integration_test.go` | Separate non-releasable public-version middleware/log/telemetry composition |
| `jsonrpc/internal_test.go` | One-excess sentinel protects duplicate bounds |
| `jsonrpc/jsonrpc.go` | Repeated raw members, per-value byte bounds and explicit metadata maps |
| `jsonrpc/jsonrpc_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `lifecycle_test.go` | Typed nils, precedence, call ordering and no partial return values |
| `log/log.go` | Deprecated slog projection forwarding without logger ownership |
| `log/log_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `multihop_integration_test.go` | Correlation preservation and per-hop causation across explicit adapters |
| `propagator.go` | Extraction before trust and generation before non-transactional injection |
| `queue/queue.go` | Application-owned maps and fresh receiving attempt identity |
| `queue/queue_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `schedule/lifecycle_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `schedule/schedule.go` | Independent roots versus explicitly propagated scheduled metadata |
| `schedule/schedule_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `telemetry/telemetry.go` | Deprecated OTel forwarding without trace/provider ownership |
| `telemetry/telemetry_test.go` | Transport/value policy, malformed input, ownership and failure assertions |
| `webhook/webhook.go` | HTTP delegation without authentication, delivery or retry |
| `webhook/webhook_test.go` | Transport/value policy, malformed input, ownership and failure assertions |

## Disposition

All 51 tracked sources and 188 exported API members (including named fields and
interface methods, aliases/constants/functions/methods) are audited, together
with seven exported private-adapter methods and private coordination contracts.
Go 1.27.2 formatting, native parsed AST/non-comment tokens including automatic
semicolons and compiler/security directives match the exact base for all files.
All 15 configured public package go-doc outputs render the clarified contracts;
the native API snapshot is their ordered stdout concatenation with one terminal
newline. Two regenerations are byte-identical. The native docs target passes
required files/local links/package rendering and all three root examples.
Other packages have no examples; unchanged non-example tests and the separately
pinned sibling harness are not rerun. Objective source-doc checks report no gaps.
Dependency/module/package records and committed vectors remain unchanged.
No runtime repair is part of this batch. Known
method-specific nil/zero behavior and trusted output versus validation distinctions
are documented rather than changed. No unresolved rationale question is currently
identified. Hosted CI, publication and any release remain coordinator-owned.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
