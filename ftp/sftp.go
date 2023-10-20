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
		return nil, errors.New("error dial is too small(min is 256 ms)")
	}
	return &sftp{cfg: cfg}, nil
}

func (c *sftp) Connect(ctx context.Context) error {
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
	defer c.close()
	err := c.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connecting to sftp: %w", err)
	}
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

func (c *sftp) close() error {
	if c.client == nil {
		return nil
	}
	if err := c.client.Close(); err != nil {
		return err
	}
	c.client = nil
	return nil
}

func (c *sftp) Rename(ctx context.Context, from, to string) error {
	defer c.close()
	err := c.Connect(ctx)
	if err != nil {
		return err
	}

	return c.client.Rename(from, to)
}
