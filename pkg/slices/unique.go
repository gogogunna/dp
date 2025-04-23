package slices

func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

func Convert[V, T any](slice []T, mapFunc func(T) V) []V {
	converted := make([]V, 0, len(slice))
	for _, v := range slice {
		converted = append(converted, mapFunc(v))
	}

	return converted
}

func UniqueFilter[T comparable](
	slice []T,
	filterFunc func(T) bool,
) []T {
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if filterFunc(v) {
			continue
		}
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}
