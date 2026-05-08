# xiter

<!-- badges -->
[![Go Version](https://img.shields.io/github/go-mod/go-version/wood-jp/xiter)](https://pkg.go.dev/github.com/wood-jp/xiter)
[![CI](https://github.com/wood-jp/xiter/actions/workflows/ci.yml/badge.svg)](https://github.com/wood-jp/xiter/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/wood-jp/xiter/badge.svg?branch=main)](https://coveralls.io/github/wood-jp/xiter?branch=main)
[![Release](https://img.shields.io/github/v/release/wood-jp/xiter)](https://github.com/wood-jp/xiter/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/wood-jp/xiter)](https://goreportcard.com/report/github.com/wood-jp/xiter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/wood-jp/xiter.svg)](https://pkg.go.dev/github.com/wood-jp/xiter)
<!-- /badges -->

Iterator utilities for `iter.Seq` and `iter.Seq2`. Wraps both sequence types with chainable methods so operations can be composed without intermediate allocations.

- [Stability](#stability)
- [Installation](#installation)
- [Usage](#usage)
  - [Wrapping sequences](#wrapping-sequences)
  - [Filter](#filter)
  - [Limit](#limit)
  - [Skip](#skip)
  - [Collect](#collect)
  - [Chaining](#chaining)
- [Contributing](#contributing)
- [Security](#security)

## Stability

v1.x releases make no breaking changes to exported APIs. New functionality may be added in minor releases; patches are bug fixes, or administrative work only.

## Installation

Go 1.26.2 or later.

```bash
go get github.com/wood-jp/xiter
```

## Usage

### Wrapping sequences

Any `iter.Seq[V]` or `iter.Seq2[K, V]` can be wrapped for chaining with a type conversion or the `ToSeq`/`ToSeq2` helpers:

```go
// type conversion
s := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5}))

// helper (useful when the compiler can't infer the type parameter)
s := xiter.ToSeq(slices.Values([]int{1, 2, 3, 4, 5}))
s2 := xiter.ToSeq2(slices.All([]string{"a", "b", "c"}))
```

Call `.Iter()` on either wrapper to get back the underlying `iter.Seq`/`iter.Seq2` for use with stdlib functions.

### Filter

`Filter` yields only the elements for which the predicate returns `true`.

```go
evens := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5})).
    Filter(func(v int) bool { return v%2 == 0 }).
    Collect()
// []int{2, 4}
```

`Seq2.Filter` receives both the key and value:

```go
for k, v := range xiter.ToSeq2(slices.All(words)).Filter(func(i int, s string) bool {
    return i%2 == 0
}) {
    fmt.Println(k, v)
}
```

### Limit

`Limit` stops iteration after at most `n` elements.

```go
first3 := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5})).
    Limit(3).
    Collect()
// []int{1, 2, 3}
```

### Skip

`Skip` discards the first `n` elements then yields the rest.

```go
after2 := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5})).
    Skip(2).
    Collect()
// []int{3, 4, 5}
```

### Collect

`Seq[V].Collect()` materializes the sequence into a `[]V`.
`Seq2[K, V].Collect()` materializes the sequence into a `map[K]V`.

```go
m := xiter.ToSeq2(slices.All([]string{"a", "b", "c"})).Collect()
// map[int]string{0: "a", 1: "b", 2: "c"}
```

### Chaining

Methods return the same wrapper type, so they compose freely:

```go
result := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8})).
    Filter(func(v int) bool { return v%2 == 0 }).
    Skip(1).
    Limit(2).
    Collect()
// []int{4, 6}
```

Early termination (`break` inside a `range` loop) propagates correctly through the chain.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

See [SECURITY.md](SECURITY.md).
