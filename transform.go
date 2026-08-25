package xiter

// Transform yields f(v) for each v in x.
func (x Seq[V]) Transform[T any](f func(V) T) Seq[T] {
	return func(yield func(T) bool) {
		x(func(v V) bool {
			return yield(f(v))
		})
	}
}

// TransformToSeq2 yields the (K, V2) pair produced by f for each v in x.
func (x Seq[V]) TransformToSeq2[K comparable, V2 any](f func(V) (K, V2)) Seq2[K, V2] {
	return func(yield func(K, V2) bool) {
		x(func(v V) bool {
			k, v2 := f(v)
			return yield(k, v2)
		})
	}
}

// Transform yields f(k, v) for each (k, v) pair in x.
func (x Seq2[K1, V1]) Transform[K2 comparable, V2 any](f func(K1, V1) (K2, V2)) Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		x(func(k K1, v V1) bool {
			return yield(f(k, v))
		})
	}
}

// TransformToSeq yields f(k, v) for each (k, v) pair in x.
func (x Seq2[K, V]) TransformToSeq[T any](f func(K, V) T) Seq[T] {
	return func(yield func(T) bool) {
		x(func(k K, v V) bool {
			return yield(f(k, v))
		})
	}
}
