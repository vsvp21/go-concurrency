package go_concurrency

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"math/rand"
	"runtime"
	"testing"
)

func TestWorkerPool_Go(t *testing.T) {
	p := NewWorkerPool(runtime.NumCPU())

	ch := make(chan int64, runtime.NumCPU())

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

	p.Go(context.TODO(), func(ctx context.Context) {
		for n := range ch {
			actualSum.Add(n)
		}
	})

	assert.Equal(t, expectedSum.Load(), actualSum.Load())
}

func BenchmarkWorkerPool_Go(b *testing.B) {
	p := NewWorkerPool(runtime.NumCPU())

	for i := 0; i < b.N; i++ {
		ch := make(chan int64, runtime.NumCPU())

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
		p.Go(context.TODO(), func(ctx context.Context) {
			for n := range ch {
				sum.Add(n)
			}
		})
	}
}
