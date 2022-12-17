package go_concurrency

import (
	"context"
	"sync"
)

func NewWorkerPool(poolSize int) *WorkerPool {
	return &WorkerPool{
		poolSize: poolSize,
		wg:       sync.WaitGroup{},
	}
}

type WorkerPool struct {
	poolSize int
	wg       sync.WaitGroup
}

func (p *WorkerPool) Go(ctx context.Context, fn func(ctx context.Context)) {
	p.wg.Add(p.poolSize)
	for i := 0; i < p.poolSize; i++ {
		go func() {
			defer p.wg.Done()
			fn(ctx)
		}()
	}

	p.wg.Wait()
}
