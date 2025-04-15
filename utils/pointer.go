package utils

// ToPointer создаёт копию значения и возвращает указатель на неё
func ToPointer[T any](v T) *T {
	return &v
}

// ToValue возвращает значение из указателя
func ToValue[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
