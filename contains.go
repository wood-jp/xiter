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
