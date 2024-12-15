package workerpool

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Basic worker pool with buffered/unbuffered channel
func workerPoolToTestAgainst(numWorkers int, tasks chan func() error) {
	var wg sync.WaitGroup
	for range numWorkers {
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
			b.Run("workerpool only exec", func(b *testing.B) {
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
				for range b.N {
					for range s.numTasks {
						if err := p.Add(task); err != nil {
							b.Fatal(err)
						}
					}
					p.Wait()
				}
			})

			// Our implementation
			b.Run("workerpool", func(b *testing.B) {
				for range b.N {
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
					for range s.numTasks {
						if err := p.Add(task); err != nil {
							b.Fatal(err)
						}
					}
					p.Wait()
				}
			})

			// Raw goroutines
			b.Run("raw_goroutines", func(b *testing.B) {
				for range b.N {
					var wg sync.WaitGroup
					var counter int32

					for range s.numTasks {
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
				for range b.N {
					var counter int32
					tasks := make(chan func() error)
					var wg sync.WaitGroup
					wg.Add(1)

					go func() {
						defer wg.Done()
						workerPoolToTestAgainst(s.poolSize, tasks)
					}()

					for range s.numTasks {
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
				for range b.N {
					var counter int32
					tasks := make(chan func() error, s.queueSize)
					var wg sync.WaitGroup
					wg.Add(1)

					go func() {
						defer wg.Done()
						workerPoolToTestAgainst(s.poolSize, tasks)
					}()

					for range s.numTasks {
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
