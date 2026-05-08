package xiter_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/wood-jp/xiter"
)

func TestSeqCollect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"empty", []int{}, nil},
		{"single", []int{42}, []int{42}},
		{"multiple", []int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[int](slices.Values(tt.input)).Collect()
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq2Collect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []string
		want  map[int]string
	}{
		{"empty", []string{}, map[int]string{}},
		{"single", []string{"a"}, map[int]string{0: "a"}},
		{"multiple", []string{"x", "y", "z"}, map[int]string{0: "x", 1: "y", 2: "z"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq2[int, string](slices.All(tt.input)).Collect()
			if !maps.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
