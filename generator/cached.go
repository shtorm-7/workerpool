package generator

func Cached[V any](scheme Scheme[V]) Scheme[V] {
	cachedValues := make([]V, 0)
	cached := false
	return func() Generator[V] {
		if !cached {
			generator := scheme()
			return func() (V, bool) {
				if value, ok := generator(); ok {
					cachedValues = append(cachedValues, value)
					return value, true
				}
				scheme = Range(cachedValues)
				cached = true
				return *new(V), false
			}
		}
		return scheme()
	}
}
