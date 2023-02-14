package generator

func Converted[V, CV any](generator Generator[V], converter func(V) CV) Generator[CV] {
	return func() (CV, bool) {
		if value, ok := generator(); ok {
			return converter(value), true
		}
		return *new(CV), false
	}
}
