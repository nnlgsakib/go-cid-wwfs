package cid

import (
	"encoding/hex"
	"strings"
	"testing"

	blake3 "lukechampine.com/blake3"
)

func TestCIDv4_GenerationAndDecode(t *testing.T) {
	data := []byte("custom v4 data")

	// Expected textual form
	sum := blake3.Sum256(data)
	hexFull := hex.EncodeToString(sum[:])
	expected := "wwfs" + hexFull[:38]

	c := NewCidV4FromBytes(data)
	if got := c.String(); got != expected {
		t.Fatalf("unexpected v4 string: got %q want %q", got, expected)
	}
	if v := c.Version(); v != 4 {
		t.Fatalf("unexpected version: got %d want 4", v)
	}

	// Decode textual form and compare
	c2, err := Decode(expected)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if !c.Equals(c2) {
		t.Fatalf("cid mismatch after decode: %v != %v", c, c2)
	}

	// Round-trip through Bytes/Cast
	out := c.Bytes()
	_, c3, err := CidFromBytes(out)
	if err != nil {
		t.Fatalf("CidFromBytes failed: %v", err)
	}
	if !c.Equals(c3) {
		t.Fatalf("cid mismatch after CidFromBytes: %v != %v", c, c3)
	}
}

func TestCIDv4_ValidationAndDigestHelpers(t *testing.T) {
	data := []byte("x")
	c := NewCidV4FromBytes(data)
	if !IsCidV4String(c.String()) {
		t.Fatalf("IsCidV4String=false for %q", c)
	}
	if v := c.Version(); v != 4 {
		t.Fatalf("unexpected version: %d", v)
	}
	if dhex, ok := CidV4DigestHex(c); !ok || len(dhex) != 38 {
		t.Fatalf("bad digest: ok=%v len=%d", ok, len(dhex))
	}

	// Uppercase acceptance
	s := c.String()
	// force uppercase hex in digest portion
	up := s[:4] + strings.ToUpper(s[4:])
	if !IsCidV4String(up) {
		t.Fatalf("uppercase digest should validate: %q", up)
	}
	c2, err := Decode(up)
	if err != nil {
		t.Fatalf("Decode uppercase failed: %v", err)
	}
	if !c.Equals(c2) {
		t.Fatalf("uppercase decode mismatch: %v != %v", c, c2)
	}
}
