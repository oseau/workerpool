package workerpool

import (
	"context"
	"sync"
)

type worker struct {
	ctx      context.Context
	tasks    chan Task
	workerWg *sync.WaitGroup
	taskWg   *sync.WaitGroup
}

func newWorker(ctx context.Context, tasks chan Task, workerWg, taskWg *sync.WaitGroup) *worker {
	return &worker{
		ctx:      ctx,
		tasks:    tasks,
		workerWg: workerWg,
		taskWg:   taskWg,
	}
}

func (w *worker) start(startWg *sync.WaitGroup) {
	go func() {
		defer w.workerWg.Done()
		startWg.Done()

		for {
			select {
			case <-w.ctx.Done():
				return
			case task := <-w.tasks:
				task.Execute()
				w.taskWg.Done()
			}
		}
	}()
}
