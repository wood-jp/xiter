package xiter_test

import (
	"context"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/wood-jp/xiter"
)

func TestFromChan_ClosedChannel(t *testing.T) {
	t.Parallel()
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	got := xiter.FromChan(t.Context(), ch).Collect()
	if want := []int{1, 2, 3}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFromChan_EarlyBreak(t *testing.T) {
	t.Parallel()
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	got := xiter.FromChan(t.Context(), ch).Limit(2).Collect()
	if want := []int{1, 2}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if v := <-ch; v != 3 {
		t.Errorf("remaining value = %d, want 3", v)
	}
}

func TestFromChan_AlreadyCancelled(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	for range 100 {
		if got := xiter.FromChan(ctx, ch).Collect(); len(got) != 0 {
			t.Fatalf("got %v, want no values", got)
		}
	}
}

func TestFromChan_CancelWhileBlocked(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		ch := make(chan int)
		done := make(chan []int)
		go func() { done <- xiter.FromChan(ctx, ch).Collect() }()
		synctest.Wait()
		cancel()
		if got := <-done; len(got) != 0 {
			t.Errorf("got %v, want no values", got)
		}
	})
}

func TestToChan_Exhausted(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		ch := xiter.Seq[int](slices.Values([]int{1, 2, 3})).ToChan(t.Context())
		var got []int
		for v := range ch {
			got = append(got, v)
		}
		if want := []int{1, 2, 3}; !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
		synctest.Wait() // goroutine has exited; channel was closed
	})
}

func TestToChan_AlreadyCancelled(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		for range 100 {
			ch := xiter.Seq[int](slices.Values([]int{1, 2, 3})).ToChan(ctx)
			if v, ok := <-ch; ok {
				t.Fatalf("received %d, want closed channel", v)
			}
		}
	})
}

func TestToChan_CancelStopsGoroutine(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		ch := xiter.Seq[int](slices.Values([]int{1, 2, 3})).ToChan(ctx)
		if v := <-ch; v != 1 {
			t.Fatalf("first value = %d, want 1", v)
		}
		cancel()
		synctest.Wait() // returns only once the goroutine has exited
		if v, ok := <-ch; ok {
			t.Errorf("received %d after cancel, want closed channel", v)
		}
	})
}

func TestChanRoundTrip(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		ctx := t.Context()
		got := xiter.FromChan(ctx, xiter.Seq[int](slices.Values([]int{1, 2, 3, 4})).ToChan(ctx)).
			Filter(func(v int) bool { return v%2 == 0 }).
			Collect()
		if want := []int{2, 4}; !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
