package constant

type (
	Pool interface {
		BaseWorker
	}

	ResizablePool interface {
		Pool

		AddWorkers(n int)
		RemoveWorkers(n int)
	}
)
