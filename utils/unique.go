package utils

func Unique[T comparable](slice []T) []T {
	keys := make(map[T]struct{})
	list := []T{}
	for _, entry := range slice {
		if _, ok := keys[entry]; !ok {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}
