package workerpool

import "context"

type Option func(*Pool)

// WithContext sets the context for the pool
func WithContext(ctx context.Context) Option {
	return func(p *Pool) {
		if ctx != nil {
			childCtx, cancel := context.WithCancel(ctx)
			p.ctx = childCtx
			p.cancel = cancel
		}
	}
}

// WithPoolSize sets the number of workers in the pool
func WithPoolSize(size int) Option {
	return func(p *Pool) {
		if size > 0 {
			p.size = size
		}
	}
}

// WithQueueSize sets the size of the task queue
func WithQueueSize(size int) Option {
	return func(p *Pool) {
		if size > 0 {
			p.queueSize = size
		}
	}
}
