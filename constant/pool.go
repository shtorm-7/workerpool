package constant

type (
	Pool interface {
		BaseWorker

		Workers() []BaseWorker
	}

	ResizablePool interface {
		Pool

		AddWorkers(n int)
		RemoveWorkers(n int)
	}
)
