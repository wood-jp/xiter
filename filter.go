package xiter

// Filter returns a Seq that yields only elements for which p returns true.
func (x Seq[V]) Filter(p func(V) bool) Seq[V] {
	return func(yield func(V) bool) {
		x(func(v V) bool {
			if !p(v) {
				return true
			}
			return yield(v)
		})
	}
}

// Filter returns a Seq2 that yields only pairs for which p returns true.
func (x Seq2[K, V]) Filter(p func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		x(func(k K, v V) bool {
			if !p(k, v) {
				return true
			}
			return yield(k, v)
		})
	}
}
