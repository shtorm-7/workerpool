package metrics

import C "github.com/shtorm-7/workerpool/constant"

func Status[T C.BaseWorker]() C.MetricHandlerFactory[T] {
	return func(worker T) C.MetricHandler {
		return func() (C.MetricName, C.Metric) {
			return "status", worker.Status().Get()
		}
	}
}

func State[T C.Worker]() C.MetricHandlerFactory[T] {
	return func(worker T) C.MetricHandler {
		return func() (C.MetricName, C.Metric) {
			return "state", worker.State().Get()
		}
	}
}

func CompletedTasks[T C.Worker]() C.MetricHandlerFactory[T] {
	return func(worker T) C.MetricHandler {
		completed := 0
		worker.State().AddCallback(
			C.Complete,
			func() {
				completed++
			},
		)
		return func() (C.MetricName, C.Metric) {
			return "completed", completed
		}
	}
}
