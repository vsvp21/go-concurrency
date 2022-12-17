package go_concurrency

import (
	"context"
)

func OrDone[T any](ctx context.Context, ch <-chan T) <-chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-ch:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return valStream
}
