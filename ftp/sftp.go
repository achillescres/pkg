package ftp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	sftplib "github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPConfig struct {
	Addr string
	cc   ssh.ClientConfig
}

type SFTP struct {
	cfg    SFTPConfig
	client *sftplib.Client
}

func NewSFTP(cfg SFTPConfig) (*SFTP, error) {
	if cfg.cc.Timeout < time.Millisecond*256 {
		return nil, errors.New("dial is too small(min is 256 ms)")
	}
	return &SFTP{cfg: cfg}, nil
}

func (c *SFTP) Connect(_ context.Context) error {
	err := c.close()
	if err != nil {
		return fmt.Errorf("close conn: %w", err)
	}

	c.client = nil
	conn, err := ssh.Dial("tcp", c.cfg.Addr, &c.cfg.cc)
	if err != nil {
		return fmt.Errorf("dialing sftp's ssh conn: %w", err)
	}

	client, err := sftplib.NewClient(conn)
	if err != nil {
		return fmt.Errorf("connecting to sftp via ssh conn: %w", err)
	}

	c.client = client
	return nil
}

func (c *SFTP) List(ctx context.Context, path string) ([]Entry, error) {
	err := c.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connecting to sftp: %w", err)
	}
	defer c.close()

	entries, err := c.client.ReadDir(path)
	owmEntries := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		owmEntries = append(owmEntries, Entry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Size:  uint64(entry.Size()),
			Time:  entry.ModTime(),
		})
	}
	return owmEntries, err
}

func (c *SFTP) UploadFile(ctx context.Context, filepath string, f io.Reader) error {
	err := c.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connecting to sftp: %w", err)
	}
	defer c.close()
	file, err := c.client.Create(filepath)
	if err != nil {
		return fmt.Errorf("couldn't create new file: %w", err)
	}
	_, err = io.Copy(file, f)
	if err != nil {
		return fmt.Errorf("couldn't write to file: %w", err)
	}

	return nil
}

func (c *SFTP) close() error {
	if c.client == nil {
		return nil
	}
	err := c.client.Close()
	if err != nil {
		return fmt.Errorf("closing sftp: %w", err)
	}
	c.client = nil
	return nil
}

func (c *SFTP) Rename(ctx context.Context, from, to string) error {
	err := c.Connect(ctx)
	if err != nil {
		return err
	}
	defer c.close()

	return c.client.Rename(from, to)
}

func (c *SFTP) DownloadFile(ctx context.Context, filepath string, f io.Writer) error {
	err := c.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connecting to sftp: %w", err)
	}
	defer c.close()

	file, err := c.client.Open(filepath)
	if err != nil {
		return fmt.Errorf("couldn't open file: %w", err)
	}
	_, err = io.Copy(f, file)
	if err != nil {
		return fmt.Errorf("couldn't copy file: %w", err)
	}

	return nil
}
