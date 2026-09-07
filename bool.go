package xiter

// All reports whether pred returns true for every element of x.
// It stops at the first false. Returns true for an empty sequence.
func (x Seq[V]) All(pred func(V) bool) bool {
	for v := range x.Iter() {
		if !pred(v) {
			return false
		}
	}
	return true
}

// Any reports whether pred returns true for at least one element of x.
// It stops at the first true. Returns false for an empty sequence.
func (x Seq[V]) Any(pred func(V) bool) bool {
	for v := range x.Iter() {
		if pred(v) {
			return true
		}
	}
	return false
}

// All reports whether pred returns true for every pair of x.
// It stops at the first false. Returns true for an empty sequence.
func (x Seq2[K, V]) All(pred func(K, V) bool) bool {
	for k, v := range x.Iter() {
		if !pred(k, v) {
			return false
		}
	}
	return true
}

// Any reports whether pred returns true for at least one pair of x.
// It stops at the first true. Returns false for an empty sequence.
func (x Seq2[K, V]) Any(pred func(K, V) bool) bool {
	for k, v := range x.Iter() {
		if pred(k, v) {
			return true
		}
	}
	return false
}
