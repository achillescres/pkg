package batch

import (
	"fmt"
)

func Use[T any](slice []T, batchSize uint, f func(batch []T) error) error {
	for start := 0; start < len(slice); start += int(batchSize) {
		err := f(slice[start : start+int(batchSize)])
		if err != nil {
			return fmt.Errorf("use batch [%d:%d]: %w", start, start+int(batchSize), err)
		}
	}
	return nil
}
