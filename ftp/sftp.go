package ftp

import (
	"context"
	"errors"
	"fmt"
	sftplib "github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"time"
)

type SFTPConfig struct {
	Addr string
	cc   ssh.ClientConfig
}

type sftp struct {
	cfg    SFTPConfig
	client *sftplib.Client
}

func NewSFTP(cfg SFTPConfig) (SFTP, error) {
	if cfg.cc.Timeout < time.Millisecond*256 {
		return nil, errors.New("dial is too small(min is 256 ms)")
	}
	return &sftp{cfg: cfg}, nil
}

func (c *sftp) Connect(_ context.Context) error {
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

func (c *sftp) UploadFile(ctx context.Context, filepath string, f io.Reader) error {
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

func (c *sftp) close() {
	if c.client == nil {
		return
	}
	c.client.Close()
	c.client = nil
	return
}

func (c *sftp) Rename(ctx context.Context, from, to string) error {
	err := c.Connect(ctx)
	if err != nil {
		return err
	}
	defer c.close()

	return c.client.Rename(from, to)
}
