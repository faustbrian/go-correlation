# Generation

The zero-option factory creates canonical UUIDv4 values through
`identifier/uuid`. Each factory owns a bounded buffer over
`crypto/rand.Reader` so concurrent request identity generation amortizes
system entropy reads without global mutable state. There is no
ambient generator and callers may inject an instance-scoped generator for
testing or a different identifier family. `NewFactory` borrows that supplied
generator for the complete lifetime of the factory. The caller must keep it
valid for that lifetime and must provide any synchronization its `New` method
requires; the factory does not copy, close, or otherwise take ownership of it.

Deterministic generation is opt-in through `NewDeterministic`. Its HMAC input
contains a package domain marker, numeric strategy version, domain length and
value, and input length and value. Keys and inputs are copied or consumed only
for the call. The checked-in vectors detect changes to this contract.

An unkeyed strategy is suitable only for public, high-entropy inputs. It makes
small input spaces enumerable and workflows linkable, so it is never a
default and must not feed metrics.
