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
