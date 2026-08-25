package xiter

import "iter"

// BoolSeq is a Seq[bool] with terminal operations And and Or, and the
// non-terminal Not. Convert any Seq[bool] with BoolSeq(s), or produce one
// directly from a Seq or Seq2 via MapBool.
type BoolSeq Seq[bool]

// Iter returns the underlying iter.Seq[bool] for interop with stdlib functions.
func (x BoolSeq) Iter() iter.Seq[bool] { return iter.Seq[bool](x) }

// And returns true if all elements are true. Returns true for an empty sequence.
func (x BoolSeq) And() bool {
	for v := range x.Iter() {
		if !v {
			return false
		}
	}
	return true
}

// Or returns true if any element is true. Returns false for an empty sequence.
func (x BoolSeq) Or() bool {
	for v := range x.Iter() {
		if v {
			return true
		}
	}
	return false
}

// Not returns a BoolSeq with each element negated.
func (x BoolSeq) Not() BoolSeq {
	return BoolSeq(Seq[bool](x).Transform(func(b bool) bool { return !b }))
}

// MapBool applies p to each element of x, returning a BoolSeq of the results.
func (x Seq[V]) MapBool(p func(V) bool) BoolSeq {
	return BoolSeq(x.Transform(p))
}

// MapBool applies p to each (k, v) pair of x, returning a BoolSeq of the results.
func (x Seq2[K, V]) MapBool(p func(K, V) bool) BoolSeq {
	return BoolSeq(x.TransformToSeq(p))
}
