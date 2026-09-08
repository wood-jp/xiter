# TODO

Ideas for additional iterator methods, grouped by what they require of the
element type.

## No constraint (`any`)

- **ForEach** — terminal consumer that runs a side-effecting func over each
  element, since not everything ends in `Collect`.
- **Reduce / Fold** — accumulate via `func(acc A, v V) A`; doesn't need `V`
  to be anything, just an accumulator and a combining func.
- **Flatten** — `Seq[Seq[V]] → Seq[V]`. Becomes load-bearing the moment
  `Windows`/`Chunks` (below) ship, since those return `Seq[Seq[V]]`.
- **Peek(f func(V))** — run a side-effecting func on each element and pass
  it through unchanged (debugging/logging taps, mirrors `io.TeeReader` /
  Rx's `Do`/`tap`). Behaviorally equivalent to
  `Transform(func(v V) V { f(v); return v })`, but `f`'s signature returns
  nothing, so the sequence is untouched by construction rather than by
  convention — makes "can't mutate what flows through" visible in the type.

## `comparable`

Following the pattern established by `Contains`.

- **Uniq / Distinct** — dedupe elements by a `key func(V) C` projection,
  yielding only first occurrences. Shares the same `key func(V) C` shape as
  `Contains`; likely the next method to build.
- **IndexOf** — like `Contains` but returns the position (or -1) of the first
  match instead of a bool.
- **CountOf** — count elements whose key equals a target (or count
  occurrences per key, returning `map[C]int`).
- **GroupBy** — partition into `map[C][]V` by key.
- **Union / Intersect / Diff** — set operations between two sequences by key,
  natural once `Uniq` and a lookup-set helper exist.

## Windowing (no constraint on `V`)

- **Windows(size int) Seq[Seq[V]]** — sliding, overlapping window of `size`
  elements per output (originally proposed as `SlidingWindowFn(size, fn)`).
  Naming follows Rust's `slice::windows` / more_itertools'
  `sliding_window`. Kotlin unifies this and `Chunks` into one
  `windowed(size, step)` call if a single primitive is preferred.
- **Chunks(size int) Seq[Seq[V]]** — non-overlapping, one output per `size`
  elements (originally proposed as `WindowFn(size, fn)`). Naming follows
  Rust's `slice::chunks` / Python's `itertools.batched`.
  Prefer keeping `fn` out of both signatures and composing with the
  existing `Transform` instead (e.g. `x.Windows(3).Transform(fn)`), to stay
  consistent with `Filter`/`Transform` being pure primitives in this repo.

  Return `Seq[Seq[V]]` rather than `Seq[[]V]`: any windower over a general
  `Seq[V]` must buffer `size` elements internally regardless of return
  type, so `Seq[[]V]` forces a choice between allocating a fresh copy per
  window (`O(n·size)` allocs) or reusing/mutating one buffer (a
  `bufio.Scanner.Bytes()`-style footgun if the caller retains a window past
  one iteration). `Seq[Seq[V]]` keeps that choice internal — each window is
  a lazy sub-sequence over the shared buffer, and range-over-func's
  synchronous execution guarantees the outer loop can't slide the window
  until the caller has drained (or broken out of) the inner `Seq[V]`, so
  reuse-without-copy is safe by construction. Callers who want a `[]V`
  snapshot call `.Collect()` on the inner seq themselves.

## Two sequences at once

File organization: these belong together in a new `combine.go` (mirroring the
`xiter.go` / `filter.go` / `transform.go` / `terminal.go` / `keyed.go` split
from the current codebase), created once the first of these methods ships —
no empty placeholder file in the meantime. Update `README.md`'s TOC and
per-func source links to match when that happens.

- **Zip** — `Seq[A], Seq[B] → Seq2[A, B]` (or a paired struct). No
  constraint needed on `A`/`B`.
- **Merge** — combine two already-sorted sequences into one sequence.
  TBD on how to execute the merge, but avoid sorting anything.
- **Equal** — compare two sequences element-by-element for equality (mirrors
  `slices.Equal`). Needs `comparable` on `V`.
