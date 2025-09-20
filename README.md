go-cid
==================

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](http://ipn.io)
[![](https://img.shields.io/badge/project-IPFS-blue.svg?style=flat-square)](http://ipfs.io/)
[![](https://img.shields.io/badge/freenode-%23ipfs-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23ipfs)
[![](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)
[![GoDoc](https://godoc.org/github.com/nnlgsakib/go-cid-wwfs?status.svg)](https://godoc.org/github.com/nnlgsakib/go-cid-wwfs)
[![Coverage Status](https://coveralls.io/repos/github/ipfs/go-cid/badge.svg?branch=master)](https://coveralls.io/github/ipfs/go-cid?branch=master)
[![Travis CI](https://travis-ci.org/ipfs/go-cid.svg?branch=master)](https://travis-ci.org/ipfs/go-cid)

> A package to handle content IDs in Go.

This is an implementation in Go of the [CID spec](https://github.com/ipld/cid).
It is used in `go-ipfs` and related packages to refer to a typed hunk of data.

## Lead Maintainer

[Eric Myhre](https://github.com/warpfork)

## Table of Contents

- [Install](#install)
- [Usage](#usage)
 - [CID v2](#cid-v2)
- [API](#api)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-cid` is a standard Go module which can be installed with:

```sh
go get github.com/nnlgsakib/go-cid-wwfs
```

## Usage

### Running tests

Run tests with `go test` from the directory root

```sh
go test
```

### Examples

#### Parsing string input from users

```go
// Create a cid from a marshaled string
c, err := cid.Decode("bafzbeigai3eoy2ccc7ybwjfz5r3rdxqrinwi4rwytly24tdbh6yk7zslrm")
if err != nil {...}

fmt.Println("Got CID: ", c)
```

#### Creating a CID from scratch

```go

import (
  cid "github.com/nnlgsakib/go-cid-wwfs"
  mc "github.com/multiformats/go-multicodec"
  mh "github.com/multiformats/go-multihash"
)

// Create a cid manually by specifying the 'prefix' parameters
pref := cid.Prefix{
	Version: 1,
	Codec: uint64(mc.Raw),
	MhType: mh.SHA2_256,
	MhLength: -1, // default length
}

// And then feed it some data
c, err := pref.Sum([]byte("Hello World!"))
if err != nil {...}

fmt.Println("Created CID: ", c)
```

#### Check if two CIDs match

```go
// To test if two cid's are equivalent, be sure to use the 'Equals' method:
if c1.Equals(c2) {
	fmt.Println("These two refer to the same exact data!")
}
```

#### Check if some data matches a given CID

```go
// To check if some data matches a given cid, 
// Get your CIDs prefix, and use that to sum the data in question:
other, err := c.Prefix().Sum(mydata)
if err != nil {...}

if !c.Equals(other) {
    fmt.Println("This data is different.")
}

```

## CID v2

This module additionally supports a custom CID version 2 textual format for specialized use-cases. CID v2 values are 42-character strings:

- Prefix: `wwfs`
- Digest: first 38 lowercase-hex characters of a BLAKE3-256 digest

Creation and parsing:

```go
import (
    cid "github.com/nnlgsakib/go-cid-wwfs"
)

data := []byte("custom v2 data")
v2 := cid.NewCidV2FromBytes(data)
fmt.Println(v2.String()) // e.g. "wwfs..." (42 chars)
fmt.Println(v2.Version()) // 2

// Parse textual form back to a CID
same, err := cid.Decode(v2.String())
if err != nil { /* handle */ }
fmt.Println(v2.Equals(same)) // true
```

Notes:
- CID v0 and v1 behavior remains unchanged.
- CID v2 does not use multibase or multicodec; `String()`, `Encode()`, and `StringOfBase()` return the canonical textual form.
- `CidFromBytes` recognizes the 42-byte ASCII form and returns CID v2.

## Contribute

PRs are welcome!

Small note: If editing the Readme, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © Jeromy Johnson
