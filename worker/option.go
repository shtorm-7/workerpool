package worker

import C "github.com/shtorm-7/workerpool/constant"

type WorkerOption func(worker *Worker)

func WithFlow(flow Flow) WorkerOption {
	return func(worker *Worker) {
		worker.flow = flow
	}
}

func WithMetrics(factories ...C.MetricHandlerFactory[C.Worker]) WorkerOption {
	return func(worker *Worker) {
		for _, factory := range factories {
			worker.metricHandlers = append(worker.metricHandlers, factory(worker))
		}
	}
}
