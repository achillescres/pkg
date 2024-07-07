package batch

import (
	"fmt"
)

func Use[T any](slice []T, batchSize uint, f func(batch []T) error) error {
	for start := 0; start < len(slice); {
		next := min(start+int(batchSize), len(slice))
		err := f(slice[start:next])
		if err != nil {
			return fmt.Errorf("use batch [%d:%d]: %w", start, start+int(batchSize), err)
		}

		// next batch
		start = next
	}
	return nil
}
