package generator

func Value[V any](value V) Scheme[V] {
	return func() Generator[V] {
		ok := true
		return func() (V, bool) {
			if ok {
				ok = false
				return value, true
			}
			return *new(V), false
		}
	}
}
