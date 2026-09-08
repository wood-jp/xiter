package xiter_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/wood-jp/xiter"
)

func TestSeqFilter_Integers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []int
		pred  func(int) bool
		want  []int
	}{
		{"all pass", []int{1, 2, 3, 4, 5}, func(v int) bool { return true }, []int{1, 2, 3, 4, 5}},
		{"none pass", []int{1, 2, 3, 4, 5}, func(v int) bool { return false }, nil},
		{"evens only", []int{1, 2, 3, 4, 5, 6}, func(v int) bool { return v%2 == 0 }, []int{2, 4, 6}},
		{"odds only", []int{1, 2, 3, 4, 5}, func(v int) bool { return v%2 != 0 }, []int{1, 3, 5}},
		{"greater than 3", []int{1, 2, 3, 4, 5}, func(v int) bool { return v > 3 }, []int{4, 5}},
		{"empty input", []int{}, func(v int) bool { return true }, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[int](slices.Values(tt.input)).Filter(tt.pred).Collect()
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeqFilter_Strings(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []string
		pred  func(string) bool
		want  []string
	}{
		{"non-empty", []string{"a", "", "b", "", "c"}, func(v string) bool { return v != "" }, []string{"a", "b", "c"}},
		{"has prefix", []string{"foo", "bar", "foobar", "baz"}, func(v string) bool { return strings.HasPrefix(v, "foo") }, []string{"foo", "foobar"}},
		{"empty input", []string{}, func(v string) bool { return true }, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[string](slices.Values(tt.input)).Filter(tt.pred).Collect()
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

type person struct {
	name string
	age  int
}

func TestSeqFilter_CustomTypes(t *testing.T) {
	t.Parallel()
	people := []person{
		{"Alice", 30},
		{"Bob", 17},
		{"Charlie", 25},
		{"Dave", 15},
	}
	tests := []struct {
		name string
		pred func(person) bool
		want []person
	}{
		{"adults", func(p person) bool { return p.age >= 18 }, []person{{"Alice", 30}, {"Charlie", 25}}},
		{"minors", func(p person) bool { return p.age < 18 }, []person{{"Bob", 17}, {"Dave", 15}}},
		{"none match", func(p person) bool { return p.age > 100 }, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xiter.Seq[person](slices.Values(people)).Filter(tt.pred).Collect()
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeqFilter_Pointers(t *testing.T) {
	t.Parallel()
	a, b, c, d := 1, 2, 3, 4
	input := []*int{&a, &b, &c, &d}

	tests := []struct {
		name string
		pred func(*int) bool
		want []int
	}{
		{"even values", func(p *int) bool { return *p%2 == 0 }, []int{2, 4}},
		{"odd values", func(p *int) bool { return *p%2 != 0 }, []int{1, 3}},
		{"all pass", func(p *int) bool { return true }, []int{1, 2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got []int
			for v := range xiter.Seq[*int](slices.Values(input)).Filter(tt.pred) {
				got = append(got, *v)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeqFilter_EarlyTermination(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     []int
		pred      func(int) bool
		stopAfter int
		want      []int
	}{
		{"stop after 2 evens", []int{1, 2, 3, 4, 5, 6}, func(v int) bool { return v%2 == 0 }, 2, []int{2, 4}},
		{"stop after 1", []int{1, 2, 3, 4, 5}, func(v int) bool { return true }, 1, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got []int
			for v := range xiter.Seq[int](slices.Values(tt.input)).Filter(tt.pred) {
				got = append(got, v)
				if len(got) == tt.stopAfter {
					break
				}
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeqFilter_StatefulPredicate(t *testing.T) {
	t.Parallel()

	t.Run("first n matching", func(t *testing.T) {
		t.Parallel()
		input := []int{1, 2, 3, 4, 5, 6, 7, 8}
		count := 0
		pred := func(v int) bool {
			if v%2 == 0 && count < 2 {
				count++
				return true
			}
			return false
		}
		got := xiter.Seq[int](slices.Values(input)).Filter(pred).Collect()
		want := []int{2, 4}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("alternating", func(t *testing.T) {
		t.Parallel()
		input := []int{1, 2, 3, 4, 5, 6}
		i := 0
		pred := func(_ int) bool {
			i++
			return i%2 == 1
		}
		got := xiter.Seq[int](slices.Values(input)).Filter(pred).Collect()
		want := []int{1, 3, 5}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestSeq2Filter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []string
		pred     func(int, string) bool
		wantKeys []int
		wantVals []string
	}{
		{
			"even indices",
			[]string{"a", "b", "c", "d", "e"},
			func(k int, _ string) bool { return k%2 == 0 },
			[]int{0, 2, 4},
			[]string{"a", "c", "e"},
		},
		{
			"filter on value",
			[]string{"foo", "bar", "foobar", "baz"},
			func(_ int, v string) bool { return strings.HasPrefix(v, "foo") },
			[]int{0, 2},
			[]string{"foo", "foobar"},
		},
		{
			"filter on both",
			[]string{"a", "b", "c", "d", "e"},
			func(k int, v string) bool { return k > 1 && v != "d" },
			[]int{2, 4},
			[]string{"c", "e"},
		},
		{
			"none pass",
			[]string{"a", "b", "c"},
			func(_ int, _ string) bool { return false },
			nil,
			nil,
		},
		{
			"all pass",
			[]string{"a", "b", "c"},
			func(_ int, _ string) bool { return true },
			[]int{0, 1, 2},
			[]string{"a", "b", "c"},
		},
		{
			"empty input",
			[]string{},
			func(_ int, _ string) bool { return true },
			nil,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var gotKeys []int
			var gotVals []string
			for k, v := range xiter.Seq2[int, string](slices.All(tt.input)).Filter(tt.pred) {
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

func TestSeq2Filter_EarlyTermination(t *testing.T) {
	t.Parallel()
	input := []string{"a", "b", "c", "d", "e", "f"}
	var gotKeys []int
	var gotVals []string
	for k, v := range xiter.Seq2[int, string](slices.All(input)).Filter(func(k int, _ string) bool { return k%2 == 0 }) {
		gotKeys = append(gotKeys, k)
		gotVals = append(gotVals, v)
		if len(gotKeys) == 2 {
			break
		}
	}
	if !slices.Equal(gotKeys, []int{0, 2}) {
		t.Errorf("keys: got %v, want %v", gotKeys, []int{0, 2})
	}
	if !slices.Equal(gotVals, []string{"a", "c"}) {
		t.Errorf("vals: got %v, want %v", gotVals, []string{"a", "c"})
	}
}

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
