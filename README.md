# xiter

<!-- badges -->
[![Go Version](https://img.shields.io/github/go-mod/go-version/wood-jp/xiter)](https://pkg.go.dev/github.com/wood-jp/xiter)
[![CI](https://github.com/wood-jp/xiter/actions/workflows/ci.yml/badge.svg)](https://github.com/wood-jp/xiter/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/wood-jp/xiter/badge.svg?branch=main)](https://coveralls.io/github/wood-jp/xiter?branch=main)
[![Release](https://img.shields.io/github/v/release/wood-jp/xiter)](https://github.com/wood-jp/xiter/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/wood-jp/xiter.svg)](https://pkg.go.dev/github.com/wood-jp/xiter)
<!-- /badges -->

Iterator utilities for `iter.Seq` and `iter.Seq2`. Wraps both sequence types with chainable methods so operations can be composed without intermediate allocations.

- [Stability](#stability)
- [Installation](#installation)
- [Usage](#usage)
  - [Wrapping sequences](#wrapping-sequences)
  - [Filter](#filter)
  - [Transform](#transform)
  - [Limit](#limit)
  - [Skip](#skip)
  - [Collect](#collect)
  - [Contains](#contains)
  - [All and Any](#all-and-any)
  - [Chaining](#chaining)
  - [Changing element type](#changing-element-type)
- [Contributing](#contributing)
- [Security](#security)

## Stability

v1.x releases make no breaking changes to exported APIs. New functionality may be added in minor releases; patches are bug fixes, or administrative work only.

## Installation

Go 1.27.0 or later.

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

### Transform

`Transform` maps each element of a `Seq[V]` through a function, yielding a `Seq[T]`.

```go
strs := xiter.Seq[int](slices.Values([]int{1, 2, 3})).
    Transform(strconv.Itoa).
    Collect()
// []string{"1", "2", "3"}
```

`Seq2.Transform` maps each `(K1, V1)` pair of a `Seq2` through a function, yielding a `Seq2[K2, V2]`:

```go
upper := xiter.ToSeq2(slices.All([]string{"a", "b", "c"})).
    Transform(func(k int, v string) (int, string) { return k, strings.ToUpper(v) })
```

`TransformToSeq2` converts a `Seq[V]` to a `Seq2[K, V2]` (e.g. pairing each element with a derived key). `TransformToSeq` does the reverse, collapsing each `(K, V)` pair into a single `T`.

```go
// pair each word with its length as the key
withLen := xiter.Seq[string](slices.Values([]string{"go", "rust", "zig"})).
    TransformToSeq2(func(v string) (int, string) { return len(v), v })
// Seq2[int, string]: (2,"go"), (4,"rust"), (3,"zig")
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

### Contains

`Contains` reports whether a sequence has an element whose key equals a target, stopping at the first match. Pass a key extractor; use the identity function when the element type is already comparable:

```go
found := xiter.Seq[int](slices.Values([]int{1, 2, 3})).
    Contains(func(v int) int { return v }, 2)
// true
```

For a type that isn't comparable, or to match on a derived key, extract that key instead:

```go
found := xiter.Seq[person](slices.Values(people)).
    Contains(func(p person) string { return p.name }, "alice")
```

`Seq2.Contains` receives both the key and value when computing the comparison key.

### All and Any

`All` reports whether a predicate holds for every element, stopping at the first `false`. `Any` reports whether it holds for at least one, stopping at the first `true`.

```go
// are all even numbers > 0?
ok := xiter.Seq[int](slices.Values([]int{2, 4, 6})).
    All(func(v int) bool { return v > 0 })
// true

// does any word start with "go"?
found := xiter.Seq[string](slices.Values([]string{"rust", "go", "zig"})).
    Any(func(v string) bool { return strings.HasPrefix(v, "go") })
// true

// no odd numbers in the slice?
noneOdd := xiter.Seq[int](slices.Values([]int{2, 4, 6})).
    All(func(v int) bool { return v%2 == 0 })
// true
```

`Seq2.All` and `Seq2.Any` receive both the key and value:

```go
allEvenIndices := xiter.ToSeq2(slices.All([]string{"a", "b", "c"})).
    All(func(k int, _ string) bool { return k%2 == 0 })
// false
```

Empty sequences: `All` returns `true`, `Any` returns `false`.

### Chaining

Methods return the same wrapper type, so they compose freely:

```go
result := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8})).
    Filter(func(v int) bool { return v%2 == 0 }).
    Skip(1).
    Limit(2).
    Transform(strconv.Itoa).
    Collect()
// []string{"4", "6"}
```

Early termination (`break` inside a `range` loop) propagates correctly through the chain.

### Changing element type

It is not possible to cast from `Seq[A]` to `Seq[B]` as these are are different underlying function types. Use `Transform` instead:

```go
anySeq := xiter.Seq[any](slices.Values([]any{"a", "b", "c"}))

strs := anySeq.Transform(func(v any) string { return v.(string) }).Collect()
// []string{"a", "b", "c"}
```

Note that in this example, if an element wasn't actually a `string` the type assertion would panic. That is entirely the choice of the caller, who must provide the transformation function.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

See [SECURITY.md](SECURITY.md).

## Attribution

*This library contains some code based on [zkr-go-common-public/iter](https://github.com/zircuit-labs/zkr-go-common-public/tree/main/iter)*
