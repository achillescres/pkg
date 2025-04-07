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
	c.close()

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

func (c *SFTP) close() {
	if c.client == nil {
		return
	}
	c.client.Close()
	c.client = nil
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
