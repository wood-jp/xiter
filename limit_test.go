package xiter_test

import (
	"slices"
	"testing"

	"github.com/wood-jp/xiter"
)

func TestSeqLimit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []int
		limit uint
		want  []int
	}{
		{"zero limit", []int{1, 2, 3}, 0, nil},
		{"limit less than len", []int{1, 2, 3, 4, 5}, 3, []int{1, 2, 3}},
		{"limit equal to len", []int{1, 2, 3}, 3, []int{1, 2, 3}},
		{"limit greater than len", []int{1, 2, 3}, 5, []int{1, 2, 3}},
		{"empty input", []int{}, 3, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[int](slices.Values(tt.input)).Limit(tt.limit).Collect()
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq2Limit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []string
		limit    uint
		wantKeys []int
		wantVals []string
	}{
		{"zero limit", []string{"a", "b", "c"}, 0, nil, nil},
		{"limit less than len", []string{"a", "b", "c", "d", "e"}, 3, []int{0, 1, 2}, []string{"a", "b", "c"}},
		{"limit equal to len", []string{"a", "b", "c"}, 3, []int{0, 1, 2}, []string{"a", "b", "c"}},
		{"limit greater than len", []string{"a", "b", "c"}, 5, []int{0, 1, 2}, []string{"a", "b", "c"}},
		{"empty input", []string{}, 3, nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var gotKeys []int
			var gotVals []string
			for k, v := range xiter.Seq2[int, string](slices.All(tt.input)).Limit(tt.limit) {
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
