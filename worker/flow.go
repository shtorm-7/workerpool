package worker

import C "github.com/shtorm-7/workerpool/constant"

type Flow func(w *Worker)

func DefaultFlow(w *Worker) {
	for {
		select {
		case <-w.status.Await(C.Stopping):
			return
		default:
		}
		select {
		case task, ok := <-w.queue:
			if !ok {
				panic("queue is closed")
			}
			w.processTask(task)
		case <-w.status.Await(C.Stopping):
			return
		}
	}
}

func GracefulFlow(w *Worker) {
	for {
		select {
		case task, ok := <-w.queue:
			if !ok {
				panic("queue is closed")
			}
			w.processTask(task)
		case <-w.status.Await(C.Stopping):
			select {
			case task, ok := <-w.queue:
				if !ok {
					panic("queue is closed")
				}
				w.processTask(task)
			default:
				return
			}
		}
	}
}
