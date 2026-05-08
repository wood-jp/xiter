// Package xiter provides iterator utilities for iter.Seq and iter.Seq2.
package xiter

import "iter"

// Seq is a chainable wrapper around iter.Seq.
// Convert any iter.Seq[V] with Seq[V](s), then chain methods like Limit.
type Seq[V any] func(yield func(V) bool)

// Seq2 is a chainable wrapper around iter.Seq2.
// K is constrained to comparable to support Collect into a map.
type Seq2[K comparable, V any] func(yield func(K, V) bool)

// ToSeq converts a standard iter.Seq[V] to Seq[V] for use when type inference is needed.
func ToSeq[V any](s iter.Seq[V]) Seq[V] {
	return Seq[V](s)
}

// ToSeq2 converts a standard iter.Seq2[K, V] to Seq2[K, V] for use when type inference is needed.
func ToSeq2[K comparable, V any](s iter.Seq2[K, V]) Seq2[K, V] {
	return Seq2[K, V](s)
}

// Iter returns the underlying iter.Seq[V] for interop with stdlib functions.
func (x Seq[V]) Iter() iter.Seq[V] { return iter.Seq[V](x) }

// Iter returns the underlying iter.Seq2[K, V] for interop with stdlib functions.
func (x Seq2[K, V]) Iter() iter.Seq2[K, V] { return iter.Seq2[K, V](x) }
