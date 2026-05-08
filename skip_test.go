package xiter_test

import (
	"slices"
	"testing"

	"github.com/wood-jp/xiter"
)

func TestSeqSkip(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []int
		skip  uint
		want  []int
	}{
		{"zero skip", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"skip less than len", []int{1, 2, 3, 4, 5}, 2, []int{3, 4, 5}},
		{"skip equal to len", []int{1, 2, 3}, 3, nil},
		{"skip greater than len", []int{1, 2, 3}, 5, nil},
		{"empty input", []int{}, 3, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[int](slices.Values(tt.input)).Skip(tt.skip).Collect()
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq2Skip(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []string
		skip     uint
		wantKeys []int
		wantVals []string
	}{
		{"zero skip", []string{"a", "b", "c"}, 0, []int{0, 1, 2}, []string{"a", "b", "c"}},
		{"skip less than len", []string{"a", "b", "c", "d", "e"}, 2, []int{2, 3, 4}, []string{"c", "d", "e"}},
		{"skip equal to len", []string{"a", "b", "c"}, 3, nil, nil},
		{"skip greater than len", []string{"a", "b", "c"}, 5, nil, nil},
		{"empty input", []string{}, 3, nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var gotKeys []int
			var gotVals []string
			for k, v := range xiter.Seq2[int, string](slices.All(tt.input)).Skip(tt.skip) {
				gotKeys = append(gotKeys, k)
				gotVals = append(gotVals, v)
			}
			if !slices.Equal(gotKeys, tt.wantKeys) {
				t.Errorf("keys: got %v, want %v", gotKeys, tt.wantKeys)
			}
			if !slices.Equal(gotVals, tt.wantVals) {
				t.Errorf("vals: got %v, want %v", gotVals, tt.wantVals)
			}
		})
	}
}
