package cid

import (
    "encoding/hex"

    blake3 "lukechampine.com/blake3"
)

// NewCidV2FromBytes computes a BLAKE3 digest of the input and returns
// a CIDv2 textual identifier in the form: "wwfs" + first 38 hex chars of the digest.
// The result is always 42 characters long.
func NewCidV2FromBytes(data []byte) Cid {
    sum := blake3.Sum256(data)
    hexFull := hex.EncodeToString(sum[:]) // 64 lowercase hex chars
    // Truncate to 38 hex characters as specified
    v2 := "wwfs" + hexFull[:38]
    return Cid{v2}
}

// V2Builder implements a simple builder that produces CIDv2 identifiers
// using the custom BLAKE3-based textual format.
type V2Builder struct{}

func (V2Builder) Sum(data []byte) (Cid, error) { return NewCidV2FromBytes(data), nil }

// GetCodec is not applicable for CIDv2 textual identifiers; returns 0.
func (V2Builder) GetCodec() uint64 { return 0 }

// WithCodec is a no-op for CIDv2 and returns the same builder.
func (v V2Builder) WithCodec(uint64) Builder { return v }

