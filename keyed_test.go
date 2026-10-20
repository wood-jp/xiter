package xiter_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/wood-jp/xiter"
)

func TestSeqContains_KeyExtractor(t *testing.T) {
	t.Parallel()
	people := []person{
		{"Alice", 30},
		{"Bob", 17},
		{"Charlie", 25},
	}
	tests := []struct {
		name   string
		target string
		want   bool
	}{
		{"present", "Bob", true},
		{"absent", "Dave", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[person](slices.Values(people)).
				Contains(func(p person) string { return p.name }, tt.target)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeqContains_IdentityKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  []int
		target int
		want   bool
	}{
		{"hit", []int{1, 2, 3}, 2, true},
		{"miss", []int{1, 2, 3}, 5, false},
		{"empty", []int{}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[int](slices.Values(tt.input)).
				Contains(func(v int) int { return v }, tt.target)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeqContains_EarlyTermination(t *testing.T) {
	t.Parallel()
	calls := 0
	f := func(v int) int {
		calls++
		return v
	}
	seq := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5})).Transform(f)
	got := seq.Contains(func(v int) int { return v }, 2)
	if !got {
		t.Errorf("got false, want true")
	}
	if calls != 2 {
		t.Errorf("f called %d times, want 2", calls)
	}
}

func TestSeqDistinct_KeyExtractor(t *testing.T) {
	t.Parallel()
	people := []person{
		{"Alice", 30},
		{"Bob", 17},
		{"Alice", 40},
	}
	got := xiter.Seq[person](slices.Values(people)).
		Distinct(func(p person) string { return p.name }).
		Collect()
	want := []person{
		{"Alice", 30},
		{"Bob", 17},
	}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSeqDistinct_IdentityKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"duplicates", []int{1, 2, 2, 3, 1}, []int{1, 2, 3}},
		{"no duplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"empty", []int{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[int](slices.Values(tt.input)).
				Distinct(func(v int) int { return v }).
				Collect()
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeqDistinct_EarlyTermination(t *testing.T) {
	t.Parallel()
	calls := 0
	f := func(v int) int {
		calls++
		return v
	}
	seq := xiter.Seq[int](slices.Values([]int{1, 2, 3})).Transform(f)
	for range seq.Distinct(func(v int) int { return v }) {
		break
	}
	if calls != 1 {
		t.Errorf("f called %d times, want 1", calls)
	}
}

func TestSeq2Distinct(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []string
		key   func(int, string) string
		want  map[int]string
	}{
		{
			"dedupe by value",
			[]string{"a", "b", "a", "c"},
			func(_ int, v string) string { return v },
			map[int]string{0: "a", 1: "b", 3: "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq2[int, string](slices.All(tt.input)).Distinct(tt.key).Collect()
			if !maps.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq2Contains(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  []string
		key    func(int, string) string
		target string
		want   bool
	}{
		{"hit on value", []string{"a", "b", "c"}, func(_ int, v string) string { return v }, "b", true},
		{"miss on value", []string{"a", "b", "c"}, func(_ int, v string) string { return v }, "z", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq2[int, string](slices.All(tt.input)).Contains(tt.key, tt.target)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
