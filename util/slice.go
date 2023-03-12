package util

func ContainsInSlice[T comparable](items []T, item T) bool {
	for _, it := range items {
		if it == item {
			return true
		}
	}
	return false
}
