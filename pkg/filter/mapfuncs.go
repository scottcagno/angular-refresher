package filter

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Vals[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

func Filter[V any](arr []V, fn func(i int, e V) bool) []V {
	var set []V
	for i, e := range arr {
		if fn(i, e) {
			set = append(set, e)
		}
	}
	return set
}
