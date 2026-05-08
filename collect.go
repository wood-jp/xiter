package xiter

import (
	"maps"
	"slices"
)

// Collect materializes the sequence into a slice.
func (x Seq[V]) Collect() []V { return slices.Collect(x.Iter()) }

// Collect materializes the sequence into a map.
func (x Seq2[K, V]) Collect() map[K]V { return maps.Collect(x.Iter()) }
