package xiter

// Transform yields f(v) for each v in x.
//
// In Go 1.27+ (golang/go#77273), this can become a chainable method:
//
//	func (x Seq[V]) Transform[T any](f func(V) T) Seq[T] {
//	    return func(yield func(T) bool) {
//	        x(func(v V) bool { return yield(f(v)) })
//	    }
//	}
func Transform[V, T any](x Seq[V], f func(V) T) Seq[T] {
	return func(yield func(T) bool) {
		x(func(v V) bool {
			return yield(f(v))
		})
	}
}

// Transform2 yields f(k, v) for each (k, v) pair in x.
//
// In Go 1.27+ (golang/go#77273), this can become a chainable method:
//
//	func (x Seq2[K1, V1]) Transform[K2 comparable, V2 any](
//	    f func(K1, V1) (K2, V2),
//	) Seq2[K2, V2] {
//	    return func(yield func(K2, V2) bool) {
//	        x(func(k K1, v V1) bool { return yield(f(k, v)) })
//	    }
//	}
func Transform2[K1 comparable, V1 any, K2 comparable, V2 any](
	x Seq2[K1, V1], f func(K1, V1) (K2, V2),
) Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		x(func(k K1, v V1) bool {
			return yield(f(k, v))
		})
	}
}

// TransformToSeq2 yields the (K, V2) pair produced by f for each v in x.
//
// In Go 1.27+ (golang/go#77273), this can become a chainable method:
//
//	func (x Seq[V]) TransformToSeq2[K comparable, V2 any](
//	    f func(V) (K, V2),
//	) Seq2[K, V2] {
//	    return func(yield func(K, V2) bool) {
//	        x(func(v V) bool {
//	            k, v2 := f(v)
//	            return yield(k, v2)
//	        })
//	    }
//	}
func TransformToSeq2[V any, K comparable, V2 any](
	x Seq[V], f func(V) (K, V2),
) Seq2[K, V2] {
	return func(yield func(K, V2) bool) {
		x(func(v V) bool {
			k, v2 := f(v)
			return yield(k, v2)
		})
	}
}

// TransformToSeq yields f(k, v) for each (k, v) pair in x.
//
// In Go 1.27+ (golang/go#77273), this can become a chainable method:
//
//	func (x Seq2[K, V]) TransformToSeq[T any](f func(K, V) T) Seq[T] {
//	    return func(yield func(T) bool) {
//	        x(func(k K, v V) bool { return yield(f(k, v)) })
//	    }
//	}
func TransformToSeq[K comparable, V, T any](
	x Seq2[K, V], f func(K, V) T,
) Seq[T] {
	return func(yield func(T) bool) {
		x(func(k K, v V) bool {
			return yield(f(k, v))
		})
	}
}
