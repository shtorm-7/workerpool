package worker

import C "github.com/shtorm-7/workerpool/constant"

type Flow func(worker *Worker)

func DefaultFlow(worker *Worker) {
	for {
		select {
		case <-worker.status.Await(C.Stopping):
			return
		default:
		}
		select {
		case task, ok := <-worker.queue:
			if !ok {
				panic("queue is closed")
			}
			worker.processTask(task)
		case <-worker.status.Await(C.Stopping):
			return
		}
	}
}

func GracefulFlow(worker *Worker) {
	for {
		select {
		case task, ok := <-worker.queue:
			if !ok {
				panic("queue is closed")
			}
			worker.processTask(task)
		case <-worker.status.Await(C.Stopping):
			select {
			case task, ok := <-worker.queue:
				if !ok {
					panic("queue is closed")
				}
				worker.processTask(task)
			default:
				return
			}
		}
	}
}
