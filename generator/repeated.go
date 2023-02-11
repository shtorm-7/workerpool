package generator

func Repeated[V any](scheme Scheme[V], n int) Scheme[V] {
	return func() Generator[V] {
		generator := scheme()
		i := 0
		return func() (V, bool) {
			for i != n {
				value, ok := generator()
				if !ok {
					i++
					if i != n {
						generator = scheme()
						continue
					} else {
						break
					}
				}
				return value, true
			}
			return *new(V), false
		}
	}
}
