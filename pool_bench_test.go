package workerpool

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Basic worker pool with unbuffered channel
func unbufferedPool(numWorkers int, tasks chan func() error) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				task()
			}
		}()
	}
	wg.Wait()
}

// Basic worker pool with buffered channel
func bufferedPool(numWorkers int, tasks chan func() error) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				task()
			}
		}()
	}
	wg.Wait()
}

func BenchmarkComparison(b *testing.B) {
	scenarios := []struct {
		name      string
		numTasks  int
		poolSize  int
		queueSize int
	}{
		{"small_load", 100, 4, 100},
		{"medium_load", 1000, 8, 1000},
		{"high_load", 10000, 16, 10000},
	}

	for _, s := range scenarios {
		b.Run(s.name, func(b *testing.B) {
			// Our implementation
			b.Run("workerpool", func(b *testing.B) {
				p := New(
					WithPoolSize(s.poolSize),
					WithQueueSize(s.queueSize),
				)
				defer p.Stop()

				var counter int32
				task := newTestTask(func() error {
					atomic.AddInt32(&counter, 1)
					return nil
				})

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					for j := 0; j < s.numTasks; j++ {
						if err := p.Add(task); err != nil {
							b.Fatal(err)
						}
					}
					p.Wait()
				}
			})

			// Raw goroutines
			b.Run("raw_goroutines", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					var wg sync.WaitGroup
					var counter int32

					for j := 0; j < s.numTasks; j++ {
						wg.Add(1)
						go func() {
							defer wg.Done()
							atomic.AddInt32(&counter, 1)
						}()
					}
					wg.Wait()
				}
			})

			// Unbuffered channel pool
			b.Run("unbuffered_pool", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					var counter int32
					tasks := make(chan func() error)
					var wg sync.WaitGroup
					wg.Add(1)

					go func() {
						defer wg.Done()
						unbufferedPool(s.poolSize, tasks)
					}()

					for j := 0; j < s.numTasks; j++ {
						tasks <- func() error {
							atomic.AddInt32(&counter, 1)
							return nil
						}
					}
					close(tasks)
					wg.Wait()
				}
			})

			// Buffered channel pool
			b.Run("buffered_pool", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					var counter int32
					tasks := make(chan func() error, s.queueSize)
					var wg sync.WaitGroup
					wg.Add(1)

					go func() {
						defer wg.Done()
						bufferedPool(s.poolSize, tasks)
					}()

					for j := 0; j < s.numTasks; j++ {
						tasks <- func() error {
							atomic.AddInt32(&counter, 1)
							return nil
						}
					}
					close(tasks)
					wg.Wait()
				}
			})
		})
	}
}
