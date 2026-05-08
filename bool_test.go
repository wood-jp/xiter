package xiter_test

import (
	"slices"
	"testing"

	"github.com/wood-jp/xiter"
)

func TestBoolSeqAnd(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []bool
		want  bool
	}{
		{"all true", []bool{true, true, true}, true},
		{"all false", []bool{false, false, false}, false},
		{"mixed, starts true", []bool{true, false, true}, false},
		{"single true", []bool{true}, true},
		{"single false", []bool{false}, false},
		{"empty", []bool{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.BoolSeq(xiter.Seq[bool](slices.Values(tt.input))).And()
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolSeqOr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []bool
		want  bool
	}{
		{"all true", []bool{true, true, true}, true},
		{"all false", []bool{false, false, false}, false},
		{"mixed, starts false", []bool{false, true, false}, true},
		{"single true", []bool{true}, true},
		{"single false", []bool{false}, false},
		{"empty", []bool{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.BoolSeq(xiter.Seq[bool](slices.Values(tt.input))).Or()
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolSeqNot(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []bool
		want  []bool
	}{
		{"all true", []bool{true, true}, []bool{false, false}},
		{"all false", []bool{false, false}, []bool{true, true}},
		{"mixed", []bool{true, false, true}, []bool{false, true, false}},
		{"empty", []bool{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := slices.Collect(xiter.BoolSeq(xiter.Seq[bool](slices.Values(tt.input))).Not().Iter())
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolSeqNot_Chain(t *testing.T) {
	t.Parallel()

	t.Run("Not then And", func(t *testing.T) {
		t.Parallel()
		// Not(all false) => all true => And = true
		got := xiter.BoolSeq(xiter.Seq[bool](slices.Values([]bool{false, false, false}))).Not().And()
		if !got {
			t.Errorf("got false, want true")
		}
	})

	t.Run("Not then Or", func(t *testing.T) {
		t.Parallel()
		// Not(all true) => all false => Or = false
		got := xiter.BoolSeq(xiter.Seq[bool](slices.Values([]bool{true, true, true}))).Not().Or()
		if got {
			t.Errorf("got true, want false")
		}
	})

	t.Run("double Not", func(t *testing.T) {
		t.Parallel()
		input := []bool{true, false, true}
		got := slices.Collect(xiter.BoolSeq(xiter.Seq[bool](slices.Values(input))).Not().Not().Iter())
		if !slices.Equal(got, input) {
			t.Errorf("got %v, want %v", got, input)
		}
	})
}

func TestBoolSeqIter(t *testing.T) {
	t.Parallel()
	input := []bool{true, false, true, true}
	got := slices.Collect(xiter.BoolSeq(xiter.Seq[bool](slices.Values(input))).Iter())
	if !slices.Equal(got, input) {
		t.Errorf("got %v, want %v", got, input)
	}
}

func TestSeqMapBool(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []int
		pred  func(int) bool
		want  bool
	}{
		{"all even, And", []int{2, 4, 6}, func(v int) bool { return v%2 == 0 }, true},
		{"has odd, And", []int{2, 3, 6}, func(v int) bool { return v%2 == 0 }, false},
		{"none even, Or", []int{1, 3, 5}, func(v int) bool { return v%2 == 0 }, false},
		{"has even, Or", []int{1, 2, 3}, func(v int) bool { return v%2 == 0 }, true},
		{"empty, And", []int{}, func(v int) bool { return v%2 == 0 }, true},
		{"empty, Or", []int{}, func(v int) bool { return v%2 == 0 }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			seq := xiter.Seq[int](slices.Values(tt.input)).MapBool(tt.pred)
			var got bool
			if tt.name[len(tt.name)-2:] == "Or" {
				got = seq.Or()
			} else {
				got = seq.And()
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeqMapBool_Chain(t *testing.T) {
	t.Parallel()

	t.Run("Filter then MapBool then And", func(t *testing.T) {
		t.Parallel()
		// keep only evens (2,4,6), then check all > 1: true
		got := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5, 6})).
			Filter(func(v int) bool { return v%2 == 0 }).
			MapBool(func(v int) bool { return v > 1 }).
			And()
		if !got {
			t.Errorf("got false, want true")
		}
	})

	t.Run("MapBool then Not then Or", func(t *testing.T) {
		t.Parallel()
		// all odd => MapBool(isEven) => all false => Not => all true => Or = true
		got := xiter.Seq[int](slices.Values([]int{1, 3, 5})).
			MapBool(func(v int) bool { return v%2 == 0 }).
			Not().
			Or()
		if !got {
			t.Errorf("got false, want true")
		}
	})
}

func TestSeq2MapBool(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []string
		pred  func(int, string) bool
		want  bool
		op    string
	}{
		{"all even indices, And", []string{"a", "b", "c"}, func(k int, _ string) bool { return k%2 == 0 }, false, "And"},
		{"any even index, Or", []string{"a", "b", "c"}, func(k int, _ string) bool { return k%2 == 0 }, true, "Or"},
		{"all non-empty values, And", []string{"a", "b", "c"}, func(_ int, v string) bool { return v != "" }, true, "And"},
		{"empty input, And", []string{}, func(_ int, _ string) bool { return true }, true, "And"},
		{"empty input, Or", []string{}, func(_ int, _ string) bool { return true }, false, "Or"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			seq := xiter.Seq2[int, string](slices.All(tt.input)).MapBool(tt.pred)
			var got bool
			if tt.op == "Or" {
				got = seq.Or()
			} else {
				got = seq.And()
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
