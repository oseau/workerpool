package workerpool

import (
	"context"
	"runtime"
	"sync"
)

type Pool struct {
	ctx       context.Context
	cancel    context.CancelFunc
	size      int
	queueSize int
	tasks     chan Task
	workerWg  sync.WaitGroup
	taskWg    sync.WaitGroup
}

// New creates a new worker pool
func New(opts ...Option) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		ctx:       ctx,
		cancel:    cancel,
		size:      runtime.NumCPU(),
		queueSize: runtime.NumCPU() * 2,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.tasks = make(chan Task, p.queueSize)
	workers := make([]*worker, p.size)

	// start workers automatically with size configured
	var startWg sync.WaitGroup
	for i := 0; i < p.size; i++ {
		workers[i] = newWorker(p.ctx, p.tasks, &p.workerWg, &p.taskWg)
		p.workerWg.Add(1)
		startWg.Add(1)
		workers[i].start(&startWg)
	}
	startWg.Wait() // Wait for all workers to start and ready to accept tasks
	return p
}

// Add adds a task to the pool, this will never block, it will return an error if the context is cancelled or the queue is full
func (p *Pool) Add(task Task) error {
	if p.ctx.Err() != nil {
		return ErrPoolClosed
	}

	p.taskWg.Add(1)
	select {
	case p.tasks <- task:
		return nil
	default:
		p.taskWg.Done()
		return ErrQueueFull
	}
}

// Stop gracefully shuts down the pool
func (p *Pool) Stop() {
	p.cancel()        // Signal shutdown via context
	p.workerWg.Wait() // Wait for workers to finish
}

// Wait blocks until all tasks are processed
func (p *Pool) Wait() {
	p.taskWg.Wait()
}
