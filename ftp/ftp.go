package ftp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	goftp "github.com/jlaffaye/ftp"
)

type FTPConfig struct {
	Host        string
	Port        uint
	Username    string
	Password    string
	DialTimeout time.Duration
}

type FTP struct {
	cfg  FTPConfig
	conn *goftp.ServerConn
}

func NewFTP(cfg FTPConfig) (*FTP, error) {
	if cfg.DialTimeout < time.Millisecond*256 {
		return nil, errors.New("error dial is too small(min is 256 ms)")
	}
	return &FTP{cfg: cfg}, nil
}

func (c *FTP) Connect(ctx context.Context) error {
	c.close()

	c.conn = nil
	tctx, cancel := context.WithTimeout(ctx, c.cfg.DialTimeout)
	defer cancel()

	conn, err := goftp.Dial(
		fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port),
		goftp.DialWithContext(tctx),
		goftp.DialWithTimeout(c.cfg.DialTimeout),
	)
	if err != nil {
		return err
	}
	err = conn.Login(c.cfg.Username, c.cfg.Password)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *FTP) UploadFile(ctx context.Context, filepath string, f io.Reader) error {
	defer c.close()
	err := c.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connect to ftp server: %w", err)
	}
	err = c.conn.Stor(filepath, f)
	if err != nil {
		return fmt.Errorf("use STOR command to ftp server: %w", err)
	}

	return nil
}

func (c *FTP) close() {
	if c.conn == nil {
		return
	}
	if err := c.conn.Quit(); err != nil {
		c.conn = nil
		//return fmt.Errorf("close connection to ftp server: %w", err)
	}
	c.conn = nil
}

func (c *FTP) Rename(ctx context.Context, from, to string) error {
	defer c.close()
	err := c.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connect to ftp server: %w", err)
	}

	return c.conn.Rename(from, to)
}

func (c *FTP) DownloadFile(ctx context.Context, filepath string, f io.Writer) error {
	defer c.close()
	err := c.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connect to ftp server: %w", err)
	}

	file, err := c.conn.Retr(filepath)
	if err != nil {
		return fmt.Errorf("use RETR command to ftp server: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return fmt.Errorf("couldn't copy file: %w", err)
	}

	return nil
}
