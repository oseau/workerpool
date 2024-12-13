package workerpool

import (
	"context"
	"runtime"
	"sync/atomic"
	"testing"
)

// testTask is a simple Task implementation for testing
type testTask struct {
	fn func() error
}

func (t *testTask) Execute() error {
	return t.fn()
}

func newTestTask(fn func() error) Task {
	return &testTask{fn: fn}
}

func TestPool(t *testing.T) {
	t.Run("new pool with defaults", func(t *testing.T) {
		p := New()
		if p == nil {
			t.Fatal("expected non-nil pool")
		}

		// Check default values
		if p.size != runtime.NumCPU() {
			t.Errorf("expected default size %d, got %d", runtime.NumCPU(), p.size)
		}

		if p.queueSize != runtime.NumCPU()*2 {
			t.Errorf("expected default queue size %d, got %d", runtime.NumCPU()*2, p.queueSize)
		}

		// Check channels are initialized
		if p.tasks == nil {
			t.Error("expected tasks channel to be initialized")
		}

		if cap(p.tasks) != p.queueSize {
			t.Errorf("expected tasks channel capacity %d, got %d", p.queueSize, cap(p.tasks))
		}
	})

	t.Run("new pool with pool size option", func(t *testing.T) {
		expectedSize := 4
		p := New(WithPoolSize(expectedSize))

		if p.size != expectedSize {
			t.Errorf("expected size %d, got %d", expectedSize, p.size)
		}
	})

	t.Run("new pool with queue size option", func(t *testing.T) {
		expectedQueueSize := 8
		p := New(WithQueueSize(expectedQueueSize))

		if p.queueSize != expectedQueueSize {
			t.Errorf("expected queue size %d, got %d", expectedQueueSize, p.queueSize)
		}

		if cap(p.tasks) != expectedQueueSize {
			t.Errorf("expected tasks channel capacity %d, got %d", expectedQueueSize, cap(p.tasks))
		}
	})

	t.Run("new pool with context option", func(t *testing.T) {
		ctx := context.Background()
		p := New(WithContext(ctx))

		if p.ctx == nil {
			t.Error("expected context to be set")
		}

		if p.cancel == nil {
			t.Error("expected cancel function to be set")
		}
	})

	t.Run("add task to a closed pool", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		p := New(WithContext(ctx))

		// Cancel the context immediately
		cancel()

		err := p.Add(newTestTask(func() error {
			return nil
		}))

		if err != ErrPoolClosed {
			t.Errorf("expected ErrPoolClosed, got %v", err)
		}
	})

	t.Run("tasks are executed", func(t *testing.T) {
		p := New(WithPoolSize(2), WithQueueSize(100))
		defer p.Stop()
		var counter int32

		// Add tasks that increment counter
		for i := 0; i < 5; i++ {
			err := p.Add(newTestTask(func() error {
				atomic.AddInt32(&counter, 1)
				return nil
			}))
			if err != nil {
				t.Errorf("failed to add task: %v", err)
			}
		}

		p.Wait()

		if atomic.LoadInt32(&counter) != 5 {
			t.Errorf("expected counter to be 5, got %d", counter)
		}
	})

	t.Run("add task to a full queue", func(t *testing.T) {
		p := New(
			WithPoolSize(1),
			WithQueueSize(1),
		)
		defer p.Stop()

		blocker := make(chan struct{})

		// First task - will be picked up by the worker
		if err := p.Add(newTestTask(func() error {
			<-blocker
			return nil
		})); err != nil {
			t.Errorf("unexpected error on first task: %v", err)
		}

		// Second task - will fill the queue
		if err := p.Add(newTestTask(func() error {
			return nil
		})); err != nil {
			t.Errorf("unexpected error on second task: %v", err)
		}

		// Third task - should fail with ErrQueueFull
		if err := p.Add(newTestTask(func() error {
			return nil
		})); err != ErrQueueFull {
			t.Errorf("expected ErrQueueFull, got %v", err)
		}

		close(blocker)
		p.Wait()
	})
}
