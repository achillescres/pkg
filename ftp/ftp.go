package ftp

import (
	"context"
	"errors"
	"fmt"
	goftp "github.com/jlaffaye/ftp"
	"io"
	"time"
)

type Config struct {
	Host        string
	Port        uint
	Username    string
	Password    string
	DialTimeout time.Duration
}

type Client interface {
	Connect(ctx context.Context) error
	Rename(ctx context.Context, from, to string) error
	//DownloadFile(filepath string) error
	UploadFile(ctx context.Context, filepath string, f io.Reader) error
}

type client struct {
	cfg  Config
	conn *goftp.ServerConn
}

func NewFTP(cfg Config) (Client, error) {
	if cfg.DialTimeout < time.Millisecond*256 {
		return nil, errors.New("error dial is too small(min is 256 ms)")
	}
	return &client{cfg: cfg}, nil
}

func (c *client) Connect(ctx context.Context) error {
	c.conn = nil
	tctx, cancel := context.WithTimeout(ctx, c.cfg.DialTimeout)
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	conn, err := goftp.Dial(fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port), goftp.DialWithContext(tctx))
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

func (c *client) UploadFile(ctx context.Context, filepath string, f io.Reader) error {
	defer c.close()
	err := c.Connect(ctx)
	if err != nil {
		return err
	}
	err = c.conn.Stor(filepath, f)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) close() error {
	if c.conn == nil {
		return nil
	}
	if err := c.conn.Quit(); err != nil {
		return err
	}
	c.conn = nil
	return nil
}

func (c *client) Rename(ctx context.Context, from, to string) error {
	defer c.close()
	err := c.Connect(ctx)
	if err != nil {
		return err
	}

	return c.conn.Rename(from, to)
}
