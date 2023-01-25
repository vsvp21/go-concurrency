package go_concurrency

import (
	"context"
	"sync"
)

func NewAwaitingPool[T any](poolSize int, ch <-chan T) *WorkerPool[T] {
	return &WorkerPool[T]{
		poolSize: poolSize,
		ch:       ch,
		wg:       sync.WaitGroup{},
		await:    true,
	}
}

func NewPool[T any](poolSize int, ch <-chan T) *WorkerPool[T] {
	return &WorkerPool[T]{
		poolSize: poolSize,
		ch:       ch,
		wg:       sync.WaitGroup{},
		await:    false,
	}
}

type WorkerPool[T any] struct {
	poolSize int
	ch       <-chan T
	wg       sync.WaitGroup
	await    bool
}

func (p *WorkerPool[T]) Go(ctx context.Context, fn func(ctx context.Context, v T)) {
	p.add()
	for i := 0; i < p.poolSize; i++ {
		go func() {
			defer p.done()
			for v := range OrDone[T](ctx, p.ch) {
				fn(ctx, v)
			}
		}()
	}
	p.wait()
}

func (p *WorkerPool[T]) add() {
	if p.await {
		p.wg.Add(p.poolSize)
	}
}

func (p *WorkerPool[T]) done() {
	if p.await {
		p.wg.Done()
	}
}

func (p *WorkerPool[T]) wait() {
	if p.await {
		p.wg.Wait()
	}
}
