package xiter_test

import (
	"slices"
	"testing"

	"github.com/wood-jp/xiter"
)

func TestSeqAll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []int
		pred  func(int) bool
		want  bool
	}{
		{"all true", []int{2, 4, 6}, func(v int) bool { return v%2 == 0 }, true},
		{"all false", []int{1, 3, 5}, func(v int) bool { return v%2 == 0 }, false},
		{"mixed", []int{2, 3, 4}, func(v int) bool { return v%2 == 0 }, false},
		{"single element, true", []int{2}, func(v int) bool { return v%2 == 0 }, true},
		{"single element, false", []int{1}, func(v int) bool { return v%2 == 0 }, false},
		{"empty", []int{}, func(v int) bool { return v%2 == 0 }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[int](slices.Values(tt.input)).All(tt.pred)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeqAny(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []int
		pred  func(int) bool
		want  bool
	}{
		{"all true", []int{2, 4, 6}, func(v int) bool { return v%2 == 0 }, true},
		{"all false", []int{1, 3, 5}, func(v int) bool { return v%2 == 0 }, false},
		{"mixed", []int{1, 2, 3}, func(v int) bool { return v%2 == 0 }, true},
		{"single element, true", []int{2}, func(v int) bool { return v%2 == 0 }, true},
		{"single element, false", []int{1}, func(v int) bool { return v%2 == 0 }, false},
		{"empty", []int{}, func(v int) bool { return v%2 == 0 }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[int](slices.Values(tt.input)).Any(tt.pred)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeqAll_EarlyTermination(t *testing.T) {
	t.Parallel()
	calls := 0
	f := func(v int) int {
		calls++
		return v
	}
	seq := xiter.Seq[int](slices.Values([]int{2, 4, 5, 6, 8})).Transform(f)
	got := seq.All(func(v int) bool { return v%2 == 0 })
	if got {
		t.Errorf("got true, want false")
	}
	if calls != 3 {
		t.Errorf("f called %d times, want 3", calls)
	}
}

func TestSeqAny_EarlyTermination(t *testing.T) {
	t.Parallel()
	calls := 0
	f := func(v int) int {
		calls++
		return v
	}
	seq := xiter.Seq[int](slices.Values([]int{1, 3, 4, 5, 7})).Transform(f)
	got := seq.Any(func(v int) bool { return v%2 == 0 })
	if !got {
		t.Errorf("got false, want true")
	}
	if calls != 3 {
		t.Errorf("f called %d times, want 3", calls)
	}
}

func TestSeqAll_Chain(t *testing.T) {
	t.Parallel()

	t.Run("Filter then All", func(t *testing.T) {
		t.Parallel()
		got := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5, 6})).
			Filter(func(v int) bool { return v%2 == 0 }).
			All(func(v int) bool { return v > 1 })
		if !got {
			t.Errorf("got false, want true")
		}
	})
}

func TestSeq2All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []string
		pred  func(int, string) bool
		want  bool
	}{
		{"all even indices", []string{"a", "b", "c"}, func(k int, _ string) bool { return k%2 == 0 }, false},
		{"all non-empty values", []string{"a", "b", "c"}, func(_ int, v string) bool { return v != "" }, true},
		{"predicate over both", []string{"a", "bb", "ccc"}, func(k int, v string) bool { return len(v) == k+1 }, true},
		{"empty", []string{}, func(_ int, _ string) bool { return true }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq2[int, string](slices.All(tt.input)).All(tt.pred)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq2Any(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []string
		pred  func(int, string) bool
		want  bool
	}{
		{"any even index", []string{"a", "b", "c"}, func(k int, _ string) bool { return k%2 == 0 }, true},
		{"any value equals", []string{"a", "b", "c"}, func(_ int, v string) bool { return v == "z" }, false},
		{"predicate over both", []string{"a", "bb", "ccc"}, func(k int, v string) bool { return len(v) == k+1 }, true},
		{"empty", []string{}, func(_ int, _ string) bool { return true }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq2[int, string](slices.All(tt.input)).Any(tt.pred)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
