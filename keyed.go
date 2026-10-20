package xiter

// Contains reports whether any element of x has a key equal to target.
// It stops at the first match.
func (x Seq[V]) Contains[C comparable](key func(V) C, target C) bool {
	for v := range x.Iter() {
		if key(v) == target {
			return true
		}
	}
	return false
}

// Contains reports whether any pair of x has a key equal to target.
// It stops at the first match.
func (x Seq2[K, V]) Contains[C comparable](key func(K, V) C, target C) bool {
	for k, v := range x.Iter() {
		if key(k, v) == target {
			return true
		}
	}
	return false
}

// Distinct returns a Seq that yields only the first element of x for each
// key, dropping later elements whose key repeats.
func (x Seq[V]) Distinct[C comparable](key func(V) C) Seq[V] {
	return func(yield func(V) bool) {
		seen := make(map[C]struct{})
		x(func(v V) bool {
			k := key(v)
			if _, ok := seen[k]; ok {
				return true
			}
			seen[k] = struct{}{}
			return yield(v)
		})
	}
}

// Distinct returns a Seq2 that yields only the first pair of x for each
// key, dropping later pairs whose key repeats.
func (x Seq2[K, V]) Distinct[C comparable](key func(K, V) C) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seen := make(map[C]struct{})
		x(func(k K, v V) bool {
			c := key(k, v)
			if _, ok := seen[c]; ok {
				return true
			}
			seen[c] = struct{}{}
			return yield(k, v)
		})
	}
}
