package pool

import (
	"sync"

	C "github.com/shtorm-7/workerpool/constant"
)

type (
	StartHandler func(workers []C.BaseWorker)
	StopHandler  func(workers []C.BaseWorker)
)

func SequentialStart(workers []C.BaseWorker) {
	for _, worker := range workers {
		worker.Start()
	}
}

func SequentialStop(workers []C.BaseWorker) {
	for _, worker := range workers {
		worker.Stop()
	}
}

func ParallelStart(workers []C.BaseWorker) {
	var wg sync.WaitGroup
	wg.Add(len(workers))
	for _, worker := range workers {
		worker := worker
		go func() {
			worker.Start()
			wg.Done()
		}()
	}
	wg.Wait()
}

func ParallelStop(workers []C.BaseWorker) {
	var wg sync.WaitGroup
	wg.Add(len(workers))
	for _, worker := range workers {
		worker := worker
		go func() {
			worker.Stop()
			wg.Done()
		}()
	}
	wg.Wait()
}
