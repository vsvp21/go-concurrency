package go_concurrency

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"testing"
)

func TestOrDone(t *testing.T) {
	ch := make(chan int64)

	ch1 := OrDone[int64](context.TODO(), ch)

	i := atomic.NewInt64(0)
	go func() {
		i.Add(<-ch1)
	}()

	ch <- 1

	assert.Equal(t, int64(1), i.Load())
}

func TestOrDoneCtxCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())

	ch := make(chan bool)

	ch1 := OrDone[bool](ctx, ch)

	i := atomic.NewBool(true)

	done := make(chan struct{})
	go func() {
		i.Store(<-ch1)
		done <- struct{}{}
	}()

	cancel()

	<-done

	assert.Equal(t, false, i.Load())
}
