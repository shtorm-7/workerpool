package worker

import C "github.com/shtorm-7/workerpool/constant"

type WorkerOption func(worker *Worker)

func WithFlow(flow Flow) WorkerOption {
	return func(worker *Worker) {
		worker.flow = flow
	}
}

func WithMeta(meta C.Meta) WorkerOption {
	return func(worker *Worker) {
		worker.meta = meta
	}
}
