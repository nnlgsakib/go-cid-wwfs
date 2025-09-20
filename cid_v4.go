package cid

import (
	"encoding/hex"
	"errors"
	"io"
	"strings"

	blake3 "lukechampine.com/blake3"
)

// NewCidV4FromBytes computes a BLAKE3 digest of the input and returns
// a CIDv4 textual identifier in the form: "wwfs" + first 38 hex chars of the digest.
// The result is always 42 characters long. Implementation avoids unnecessary allocations.
func NewCidV4FromBytes(data []byte) Cid {
	sum := blake3.Sum256(data)
	// Only hex-encode the first 19 bytes (38 hex chars required)
	var out [42]byte
	copy(out[:4], []byte{'w', 'w', 'f', 's'})
	hex.Encode(out[4:], sum[:19])
	return Cid{string(out[:])}
}

// V4Builder implements a simple builder that produces CIDv4 identifiers
// using the custom BLAKE3-based textual format.
type V4Builder struct{}

func (V4Builder) Sum(data []byte) (Cid, error) { return NewCidV4FromBytes(data), nil }

// GetCodec is not applicable for CIDv4 textual identifiers; returns 0.
func (V4Builder) GetCodec() uint64 { return 0 }

// WithCodec is a no-op for CIDv4 and returns the same builder.
func (v V4Builder) WithCodec(uint64) Builder { return v }

// CidV4FromReader reads exactly 42 bytes from r and returns a CIDv4
// if the bytes match the CIDv4 textual format; otherwise returns an error.
func CidV4FromReader(r io.Reader) (Cid, error) {
	var buf [42]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return Undef, ErrInvalidCid{err}
	}
	s := string(buf[:])
	if !isCidV4String(s) {
		return Undef, ErrInvalidCid{errors.New("invalid cid v4 format")}
	}
	return Cid{normalizeCidV4(s)}, nil
}

// IsCidV4String reports whether s is a valid CIDv4 textual identifier.
func IsCidV4String(s string) bool { return isCidV4String(s) }

// CidV4DigestHex returns the 38-hex-character digest portion of a CIDv4 and true.
// If c is not a CIDv4, returns ("", false).
func CidV4DigestHex(c Cid) (string, bool) {
	if c.Version() != 4 {
		return "", false
	}
	return c.str[4:42], true
}

// CanonicalizeCidV4 validates and returns the canonical lowercase form
// of a CIDv4 string. Returns an error if the input is not a valid CIDv4.
func CanonicalizeCidV4(s string) (string, error) {
	if !isCidV4String(s) {
		return "", ErrInvalidCid{errors.New("invalid cid v4 format")}
	}
	return normalizeCidV4(s), nil
}

// isCidV4String checks if the string is a valid CIDv4 format
func isCidV4String(s string) bool {
	if len(s) != 42 {
		return false
	}
	if s[:4] != "wwfs" {
		return false
	}
	// Check if the remaining 38 characters are valid hex
	for i := 4; i < 42; i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// normalizeCidV4 converts a CIDv4 string to canonical lowercase form
func normalizeCidV4(s string) string {
	if len(s) != 42 || s[:4] != "wwfs" {
		return s // Return as-is if invalid
	}
	// Convert hex portion to lowercase
	return "wwfs" + strings.ToLower(s[4:])
}
