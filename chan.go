package xiter

import "context"

// FromChan returns a Seq that yields values received from ch.
//
// The sequence ends when ch is closed, when ctx is cancelled, or when the
// consumer stops iterating. Cancellation takes precedence: once ctx is done,
// no further values are yielded even if ch has values ready. Cancellation
// ends the sequence silently, so callers who need to tell it apart from ch
// closing should check ctx.Err() afterwards.
//
// There is no Seq2 equivalent, as a channel carries a single value per send.
func FromChan[V any](ctx context.Context, ch <-chan V) Seq[V] {
	return func(yield func(V) bool) {
		for {
			if ctx.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case v, ok := <-ch:
				if !ok || !yield(v) {
					return
				}
			}
		}
	}
}

// ToChan starts a goroutine that sends each element of x on the returned
// unbuffered channel. The goroutine starts immediately, not on first receive.
//
// The channel is closed when x is exhausted or ctx is cancelled. Cancellation
// takes precedence: once ctx is done, no further values are sent even if a
// receiver is ready. The context is only checked when sending, so a source
// that blocks while producing its next element must observe ctx itself, for
// example by being built with FromChan(ctx, ...).
//
// The caller must either drain the channel or cancel ctx; otherwise the
// goroutine leaks. A panic in x is not recoverable by the caller and crashes
// the program.
//
// There is no Seq2 equivalent, as a channel carries a single value per send.
func (x Seq[V]) ToChan(ctx context.Context) <-chan V {
	ch := make(chan V)
	go func() {
		defer close(ch)
		for v := range x {
			if ctx.Err() != nil {
				return
			}
			select {
			case ch <- v:
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}
