package lox

func SliceToSlice[A any, B any](in []A, f func(item A) B) []B {
	result := make([]B, len(in))
	for i, item := range in {
		result[i] = f(item)
	}
	return result
}

// SliceRemove https://stackoverflow.com/a/37335777
func SliceRemove[T any](in []T, f func(item T) bool) ([]T, T, bool) {
	for i, item := range in {
		item := item
		if f(item) {
			return append(in[:i], in[i+1:]...), item, true
		}
	}

	return in, *new(T), false
}
