package go_concurrency

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestAwaitingWorkerPool_Go(t *testing.T) {
	ch := make(chan int64, runtime.NumCPU())
	p := NewAwaitingPool[int64](runtime.NumCPU(), ch)

	expectedSum := atomic.NewInt64(0)
	actualSum := atomic.NewInt64(0)

	nums := func() []int64 {
		nums := make([]int64, 10)
		for i := 0; i < len(nums); i++ {
			nums[i] = rand.Int63()
		}

		return nums
	}()

	go func() {
		defer close(ch)
		for _, n := range nums {
			expectedSum.Add(n)
			ch <- n
		}
	}()

	p.Go(context.TODO(), func(ctx context.Context, n int64) {
		actualSum.Add(n)
	})

	assert.Equal(t, expectedSum.Load(), actualSum.Load())
}

func BenchmarkAwaitingWorkerPool_Go(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := make(chan int64, runtime.NumCPU())
		p := NewAwaitingPool[int64](runtime.NumCPU(), ch)

		nums := func() []int64 {
			nums := make([]int64, 10)
			for i := 0; i < len(nums); i++ {
				nums[i] = rand.Int63()
			}

			return nums
		}()

		go func() {
			defer close(ch)
			for _, n := range nums {
				ch <- n
			}
		}()

		sum := atomic.NewInt64(0)
		p.Go(context.TODO(), func(ctx context.Context, n int64) {
			sum.Add(n)
		})
	}
}

func TestWorkerPool_Go(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond)
	defer cancel()

	ch := make(chan int, runtime.NumCPU())
	p := NewPool[int](runtime.NumCPU(), ch)

	actualSum := atomic.NewInt64(0)

	nums := func() []int {
		nums := make([]int, 10)
		for i := 0; i < 10; i++ {
			nums[i] = i
		}

		return nums
	}()

	go func() {
		defer close(ch)
		for _, n := range nums {
			ch <- n
		}
	}()

	p.Go(context.TODO(), func(ctx context.Context, n int) {
		actualSum.Add(int64(n))
	})

	<-ctx.Done()

	assert.Equal(t, int64(45), actualSum.Load())
}
