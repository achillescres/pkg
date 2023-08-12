package s3

import (
	"context"
)

type FilterBucket interface {
	Bucket
	GetAllFileHeadsByFilter(ctx context.Context, prefix string, filter FilterByFileHeadFunc) ([]*FileHead, error)
}

type taisBucket struct {
	Bucket
}

func NewTimeBucket(b Bucket) FilterBucket {
	return &taisBucket{b}
}

func (t *taisBucket) GetAllFileHeadsByFilter(
	ctx context.Context,
	prefix string,
	filter FilterByFileHeadFunc,
) ([]*FileHead, error) {
	fhs, err := t.Bucket.GetAllFileHeads(ctx, prefix)
	if err != nil {
		return nil, err
	}
	filtered := make([]*FileHead, 0, 8)
	for _, fh := range fhs {
		if filter(fh) {
			filtered = append(filtered, fh)
		}
	}

	return filtered, nil
}
