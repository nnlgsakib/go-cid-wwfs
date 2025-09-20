package cid

import (
	"encoding/hex"
	"testing"

	blake3 "lukechampine.com/blake3"
)

func TestCIDv2_GenerationAndDecode(t *testing.T) {
	data := []byte("custom v2 data")

	// Expected textual form
	sum := blake3.Sum256(data)
	hexFull := hex.EncodeToString(sum[:])
	expected := "wwfs" + hexFull[:38]

	c := NewCidV2FromBytes(data)
	if got := c.String(); got != expected {
		t.Fatalf("unexpected v2 string: got %q want %q", got, expected)
	}
	if v := c.Version(); v != 2 {
		t.Fatalf("unexpected version: got %d want 2", v)
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
