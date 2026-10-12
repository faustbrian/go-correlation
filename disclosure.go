package correlation

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

// DisclosureMode controls identifier disclosure to logs and telemetry.
type DisclosureMode uint8

const (
	// RedactDisclosure is the safe default and emits only a marker.
	RedactDisclosure DisclosureMode = iota
	// HashDisclosure emits a keyed, domain-separated stable token.
	HashDisclosure
	// ExposeDisclosure emits raw identifier text unchanged; Disclose itself does
	// not validate it. Use only for explicitly approved diagnostic data.
	ExposeDisclosure
)

// ErrInvalidDisclosure reports unsafe disclosure configuration.
var ErrInvalidDisclosure = errors.New("correlation: invalid disclosure policy")

// DisclosurePolicy must explicitly opt into linkable or raw output.
type DisclosurePolicy struct {
	// Mode selects redaction (zero), keyed hashing or explicit raw exposure.
	Mode DisclosureMode
	// Key is borrowed during Disclose, never copied/erased; keep it stable for
	// concurrent calls. Hash mode requires nonempty key, and every mode rejects
	// more than 1024 key bytes. This is separate from deterministic workflow keys.
	Key []byte
}

// Disclose renders an identifier according to an explicit observability
// policy. label domain-separates correlation, request, and causation hashes;
// use fixed distinct labels, not attacker-controlled label framing. Hash mode
// emits 16 HMAC-SHA-256 digest bytes as unpadded base64url (22 characters).
// Neither label nor value is validated or size-limited here; ExposeDisclosure
// returns value unchanged. After the key-size check, empty value returns empty,
// nil even for an unknown mode or missing hash key. Invalid nonempty policies
// return ErrInvalidDisclosure. Callers bound inputs and synchronize borrowed keys.
func Disclose(label, value string, policy DisclosurePolicy) (string, error) {
	if len(policy.Key) > 1024 {
		return "", ErrInvalidDisclosure
	}
	if value == "" {
		return "", nil
	}
	switch policy.Mode {
	case RedactDisclosure:
		return "[redacted]", nil
	case HashDisclosure:
		if len(policy.Key) == 0 {
			return "", fmt.Errorf("%w: hashing requires a key", ErrInvalidDisclosure)
		}
		hash := hmac.New(sha256.New, policy.Key)
		_, _ = hash.Write([]byte("go-correlation/disclosure\x00"))
		_, _ = hash.Write([]byte(label))
		_, _ = hash.Write([]byte{0})
		_, _ = hash.Write([]byte(value))
		return base64.RawURLEncoding.EncodeToString(hash.Sum(nil)[:16]), nil
	case ExposeDisclosure:
		return value, nil
	default:
		return "", ErrInvalidDisclosure
	}
}
