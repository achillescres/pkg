package ftp

import (
	"context"
	"io"
	"time"
)

type Client interface {
	Connect(ctx context.Context) error
	Rename(ctx context.Context, from, to string) error
	Uploader
	Downloader
	Lister
}

type Uploader interface {
	UploadFile(ctx context.Context, filepath string, f io.Reader) error
}

type Downloader interface {
	DownloadFile(ctx context.Context, filepath string, f io.Writer) error
}

const RelativePath = "./"

type Entry struct {
	Name  string
	IsDir bool
	Size  uint64
	Time  time.Time
}

type Lister interface {
	List(ctx context.Context, path string) ([]Entry, error)
}
