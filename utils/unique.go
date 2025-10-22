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

func UniqueStringKey[T any](slice []T, key func(T) string) []T {
	keys := make(map[string]struct{})
	var list []T
	for _, entry := range slice {
		if _, ok := keys[key(entry)]; !ok {
			keys[key(entry)] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}
