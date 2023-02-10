package generator

func Converted[V, CV any](scheme Scheme[V], converter func(V) CV) Scheme[CV] {
	return func() Generator[CV] {
		generator := scheme()
		return func() (CV, bool) {
			if value, ok := generator(); ok {
				return converter(value), true
			}
			return *new(CV), false
		}
	}
}
