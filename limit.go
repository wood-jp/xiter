package xiter

// Limit returns a Seq that yields at most l elements then stops iteration.
func (x Seq[V]) Limit(l uint) Seq[V] {
	return func(yield func(V) bool) {
		var count uint
		x(func(v V) bool {
			if count >= l {
				return false
			}
			count++
			return yield(v)
		})
	}
}

// Limit returns a Seq2 that yields at most l elements then stops iteration.
func (x Seq2[K, V]) Limit(l uint) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var count uint
		x(func(k K, v V) bool {
			if count >= l {
				return false
			}
			count++
			return yield(k, v)
		})
	}
}
