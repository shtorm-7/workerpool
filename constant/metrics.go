package constant

type (
	MetricName string
	Metric     any

	Metrics map[MetricName]Metric

	MetricHandler func() (MetricName, Metric)

	MetricHandlerFactory[T BaseWorker] func(worker T) MetricHandler
)
