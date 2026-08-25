package xiter_test

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/wood-jp/xiter"
)

type personDTO struct {
	displayName string
	adult       bool
}

func TestTransform_IntToString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []int
		want  []string
	}{
		{"normal", []int{1, 2, 3}, []string{"1", "2", "3"}},
		{"single element", []int{42}, []string{"42"}},
		{"empty input", []int{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[int](slices.Values(tt.input)).Transform(strconv.Itoa).Collect()
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransform_StructToStruct(t *testing.T) {
	t.Parallel()
	people := []person{
		{"Alice", 30},
		{"Bob", 17},
		{"Charlie", 25},
	}
	toDTO := func(p person) personDTO {
		return personDTO{displayName: strings.ToUpper(p.name), adult: p.age >= 18}
	}
	got := xiter.Seq[person](slices.Values(people)).Transform(toDTO).Collect()
	want := []personDTO{
		{"ALICE", true},
		{"BOB", false},
		{"CHARLIE", true},
	}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTransform_EarlyTermination(t *testing.T) {
	t.Parallel()
	calls := 0
	f := func(v int) string {
		calls++
		return strconv.Itoa(v)
	}
	seq := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5})).Transform(f)
	var got []string
	for v := range seq {
		got = append(got, v)
		if len(got) == 2 {
			break
		}
	}
	if !slices.Equal(got, []string{"1", "2"}) {
		t.Errorf("got %v, want [1 2]", got)
	}
	if calls != 2 {
		t.Errorf("f called %d times, want 2", calls)
	}
}

func TestTransform_ChainsWithFilter(t *testing.T) {
	t.Parallel()
	seq := xiter.Seq[int](slices.Values([]int{1, 2, 3, 4, 5})).
		Filter(func(v int) bool { return v%2 == 0 })
	got := seq.Transform(strconv.Itoa).Collect()
	want := []string{"2", "4"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSeq2Transform_KeyAndValue(t *testing.T) {
	t.Parallel()
	input := []string{"a", "b", "c"}
	f := func(k int, v string) (string, int) {
		return v, k * 10
	}
	var gotKeys []string
	var gotVals []int
	for k, v := range xiter.Seq2[int, string](slices.All(input)).Transform(f) {
		gotKeys = append(gotKeys, k)
		gotVals = append(gotVals, v)
	}
	if !slices.Equal(gotKeys, []string{"a", "b", "c"}) {
		t.Errorf("keys: got %v, want [a b c]", gotKeys)
	}
	if !slices.Equal(gotVals, []int{0, 10, 20}) {
		t.Errorf("vals: got %v, want [0 10 20]", gotVals)
	}
}

func TestSeq2Transform_ValueOnly(t *testing.T) {
	t.Parallel()
	input := []string{"hello", "world"}
	f := func(k int, v string) (int, string) {
		return k, strings.ToUpper(v)
	}
	var gotKeys []int
	var gotVals []string
	for k, v := range xiter.Seq2[int, string](slices.All(input)).Transform(f) {
		gotKeys = append(gotKeys, k)
		gotVals = append(gotVals, v)
	}
	if !slices.Equal(gotKeys, []int{0, 1}) {
		t.Errorf("keys: got %v, want [0 1]", gotKeys)
	}
	if !slices.Equal(gotVals, []string{"HELLO", "WORLD"}) {
		t.Errorf("vals: got %v, want [HELLO WORLD]", gotVals)
	}
}

func TestSeq2Transform_EarlyTermination(t *testing.T) {
	t.Parallel()
	calls := 0
	f := func(k int, v string) (int, string) {
		calls++
		return k, strings.ToUpper(v)
	}
	seq := xiter.Seq2[int, string](slices.All([]string{"a", "b", "c", "d", "e"})).Transform(f)
	var gotKeys []int
	var gotVals []string
	for k, v := range seq {
		gotKeys = append(gotKeys, k)
		gotVals = append(gotVals, v)
		if len(gotKeys) == 2 {
			break
		}
	}
	if !slices.Equal(gotKeys, []int{0, 1}) {
		t.Errorf("keys: got %v, want [0 1]", gotKeys)
	}
	if !slices.Equal(gotVals, []string{"A", "B"}) {
		t.Errorf("vals: got %v, want [A B]", gotVals)
	}
	if calls != 2 {
		t.Errorf("f called %d times, want 2", calls)
	}
}

func TestSeq2Transform_ChainsWithFilter(t *testing.T) {
	t.Parallel()
	input := []string{"a", "b", "c", "d", "e"}
	seq := xiter.Seq2[int, string](slices.All(input)).
		Filter(func(k int, _ string) bool { return k%2 == 0 })
	f := func(k int, v string) (int, string) {
		return k, strings.ToUpper(v)
	}
	var gotKeys []int
	var gotVals []string
	for k, v := range seq.Transform(f) {
		gotKeys = append(gotKeys, k)
		gotVals = append(gotVals, v)
	}
	if !slices.Equal(gotKeys, []int{0, 2, 4}) {
		t.Errorf("keys: got %v, want [0 2 4]", gotKeys)
	}
	if !slices.Equal(gotVals, []string{"A", "C", "E"}) {
		t.Errorf("vals: got %v, want [A C E]", gotVals)
	}
}

func TestTransformToSeq2_Basic(t *testing.T) {
	t.Parallel()
	input := []string{"a", "b", "c"}
	seq := xiter.Seq[string](slices.Values(input)).
		TransformToSeq2(func(v string) (int, string) { return len(v), v })
	var gotKeys []int
	var gotVals []string
	for k, v := range seq {
		gotKeys = append(gotKeys, k)
		gotVals = append(gotVals, v)
	}
	if !slices.Equal(gotKeys, []int{1, 1, 1}) {
		t.Errorf("keys: got %v, want [1 1 1]", gotKeys)
	}
	if !slices.Equal(gotVals, []string{"a", "b", "c"}) {
		t.Errorf("vals: got %v, want [a b c]", gotVals)
	}
}

func TestTransformToSeq2_EarlyTermination(t *testing.T) {
	t.Parallel()
	calls := 0
	f := func(v string) (int, string) {
		calls++
		return len(v), v
	}
	seq := xiter.Seq[string](slices.Values([]string{"a", "b", "c", "d", "e"})).
		TransformToSeq2(f)
	var gotKeys []int
	for k := range seq {
		gotKeys = append(gotKeys, k)
		if len(gotKeys) == 2 {
			break
		}
	}
	if calls != 2 {
		t.Errorf("f called %d times, want 2", calls)
	}
}

func TestTransformToSeq_Basic(t *testing.T) {
	t.Parallel()
	seq := xiter.Seq2[int, string](slices.All([]string{"x", "y", "z"})).
		TransformToSeq(func(k int, v string) string { return fmt.Sprintf("%d=%s", k, v) })
	got := seq.Collect()
	want := []string{"0=x", "1=y", "2=z"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTransformToSeq_EarlyTermination(t *testing.T) {
	t.Parallel()
	calls := 0
	f := func(k int, v string) string {
		calls++
		return fmt.Sprintf("%d=%s", k, v)
	}
	seq := xiter.Seq2[int, string](slices.All([]string{"x", "y", "z", "w", "q"})).
		TransformToSeq(f)
	var got []string
	for v := range seq {
		got = append(got, v)
		if len(got) == 2 {
			break
		}
	}
	if calls != 2 {
		t.Errorf("f called %d times, want 2", calls)
	}
}

func TestTransformToSeq2_RoundTrip(t *testing.T) {
	t.Parallel()
	input := []string{"x", "y", "z"}
	got := xiter.Seq[string](slices.Values(input)).
		TransformToSeq2(func(v string) (int, string) { return len(v), strings.ToUpper(v) }).
		TransformToSeq(func(k int, v string) string { return fmt.Sprintf("%d:%s", k, v) }).
		Collect()
	want := []string{"1:X", "1:Y", "1:Z"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
