package ftp

import (
	"context"
	"io"
)

type Client interface {
	Connect(ctx context.Context) error
	Rename(ctx context.Context, from, to string) error
	Uploader
	Downloader
}

type Uploader interface {
	UploadFile(ctx context.Context, filepath string, f io.Reader) error
}

type Downloader interface {
	DownloadFile(ctx context.Context, filepath string, f io.Writer) error
}
