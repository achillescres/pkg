package batch

import (
	"testing"
)

func TestUse(t *testing.T) {
	var length int = 100

	a := make([]int, length)
	for i := 0; i < len(a); i += 1 {
		a[i] = i
	}

	batchNo := 0
	batchSize := 65000
	_ = Use(a, uint(batchSize), func(batch []int) error {
		for i := range batch {
			batch[i] = batchNo
		}
		batchNo += 1
		return nil
	})

	for i := range a {
		if a[i] != i/batchSize {
			t.Errorf("a[i] != related batchNo: %d != %d", a[i], i/batchSize)
			return
		}
	}
}
