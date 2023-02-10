package pool

import C "github.com/shtorm-7/workerpool/constant"

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

func ConcurrentStart(workers []C.BaseWorker) {
	for _, worker := range workers {
		go worker.Start()
	}
	for _, worker := range workers {
		<-worker.Status().Await(C.Started)
	}
}

func ConcurrentStop(workers []C.BaseWorker) {
	for _, worker := range workers {
		go worker.Stop()
	}
	for _, worker := range workers {
		<-worker.Status().Await(C.Stopped)
	}
}
