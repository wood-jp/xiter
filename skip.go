package xiter

// Skip returns a Seq that skips the first s elements then yields the rest.
func (x Seq[V]) Skip(s uint) Seq[V] {
	return func(yield func(V) bool) {
		var count uint
		x(func(v V) bool {
			if count < s {
				count++
				return true
			}
			return yield(v)
		})
	}
}

// Skip returns a Seq2 that skips the first s elements then yields the rest.
func (x Seq2[K, V]) Skip(s uint) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var count uint
		x(func(k K, v V) bool {
			if count < s {
				count++
				return true
			}
			return yield(k, v)
		})
	}
}
