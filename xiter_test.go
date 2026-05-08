package xiter_test

import (
	"slices"
	"testing"

	"github.com/wood-jp/xiter"
)

func TestToSeq(t *testing.T) {
	t.Parallel()
	input := []int{1, 2, 3}
	got := xiter.ToSeq(slices.Values(input)).Collect()
	if !slices.Equal(got, input) {
		t.Errorf("got %v, want %v", got, input)
	}
}

func TestToSeq2(t *testing.T) {
	t.Parallel()
	input := []string{"a", "b", "c"}
	var gotKeys []int
	var gotVals []string
	for k, v := range xiter.ToSeq2(slices.All(input)) {
		gotKeys = append(gotKeys, k)
		gotVals = append(gotVals, v)
	}
	if !slices.Equal(gotKeys, []int{0, 1, 2}) {
		t.Errorf("keys: got %v, want %v", gotKeys, []int{0, 1, 2})
	}
	if !slices.Equal(gotVals, input) {
		t.Errorf("vals: got %v, want %v", gotVals, input)
	}
}
