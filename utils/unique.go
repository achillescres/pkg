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

func UniqueFuncStringKey[T any](slice []T, key func(T) string) []T {
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

func Set[T comparable](slice []T) map[T]struct{} {
	set := make(map[T]struct{})
	for _, entry := range slice {
		set[entry] = struct{}{}
	}
	return set
}

func SetFuncStringKey[T any](slice []T, key func(T) string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, entry := range slice {
		set[key(entry)] = struct{}{}
	}
	return set
}
