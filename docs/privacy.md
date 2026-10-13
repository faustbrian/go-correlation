# Privacy

[Documentation index](README.md) · [Repository entry](../README.md)

Random correlation IDs can still connect events inside their lifetime.
Deterministic IDs create longer-lived linkability and may reveal a small input
space by enumeration. Use keyed derivation, versioned domains, narrow
retention, and access controls when stable correlation is justified.

Do not put customer identifiers, order numbers, email addresses, tenant IDs,
or other business values directly into these carriers. Prefer redaction for
logs and spans. Rotate disclosure keys independently from deterministic
workflow keys and treat both as secrets outside source control.

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
