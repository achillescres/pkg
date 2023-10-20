package ftp

import (
	"context"
	"io"
)

type FTP interface {
	Client
}

type SFTP interface {
	Client
}

type Client interface {
	Connect(ctx context.Context) error
	Rename(ctx context.Context, from, to string) error
	//DownloadFile(filepath string) error
	UploadFile(ctx context.Context, filepath string, f io.Reader) error
}
