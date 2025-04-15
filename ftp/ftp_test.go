package ftp

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockServer имитирует FTP сервер
type mockServer struct {
	listener net.Listener
	delay    time.Duration
}

func newMockServer(delay time.Duration) (*mockServer, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	s := &mockServer{
		listener: l,
		delay:    delay,
	}
	go s.serve()
	return s, nil
}

func (s *mockServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go func(c net.Conn) {
			time.Sleep(s.delay)
			c.Close()
		}(conn)
	}
}

func (s *mockServer) close() {
	s.listener.Close()
}

func (s *mockServer) addr() string {
	return s.listener.Addr().String()
}

func TestNewFTP(t *testing.T) {
	tests := []struct {
		name        string
		cfg         FTPConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			cfg: FTPConfig{
				Host:        "localhost",
				Port:        21,
				Username:    "user",
				Password:    "pass",
				DialTimeout: time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid timeout",
			cfg: FTPConfig{
				Host:        "localhost",
				Port:        21,
				Username:    "user",
				Password:    "pass",
				DialTimeout: time.Millisecond * 100,
			},
			wantErr:     true,
			errContains: "error dial is too small",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFTP(tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestFTP_UploadFile(t *testing.T) {
	// Создаем временный файл для тестирования
	content := "test content"
	tmpfile, err := os.CreateTemp("", "example")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(content))
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	tests := []struct {
		name        string
		setupFTP    func() (*FTP, error)
		filepath    string
		content     string
		wantErr     bool
		errContains string
	}{
		{
			name: "invalid credentials",
			setupFTP: func() (*FTP, error) {
				return NewFTP(FTPConfig{
					Host:        "localhost",
					Port:        2121,
					Username:    "invalid",
					Password:    "invalid",
					DialTimeout: time.Second,
				})
			},
			filepath:    "test.txt",
			content:     content,
			wantErr:     true,
			errContains: "connect to ftp server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.setupFTP()
			require.NoError(t, err)

			err = client.UploadFile(context.Background(), tt.filepath, strings.NewReader(tt.content))
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestFTP_Rename(t *testing.T) {
	tests := []struct {
		name        string
		setupFTP    func() (*FTP, error)
		from        string
		to          string
		wantErr     bool
		errContains string
	}{
		{
			name: "invalid credentials",
			setupFTP: func() (*FTP, error) {
				return NewFTP(FTPConfig{
					Host:        "localhost",
					Port:        2121,
					Username:    "invalid",
					Password:    "invalid",
					DialTimeout: time.Second,
				})
			},
			from:        "old.txt",
			to:          "new.txt",
			wantErr:     true,
			errContains: "connect to ftp server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.setupFTP()
			require.NoError(t, err)

			err = client.Rename(context.Background(), tt.from, tt.to)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestFTP_ConnectionFailures(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*FTP, func(), error)
		wantErr     bool
		errContains string
	}{
		{
			name: "server shutdown immediately",
			setup: func() (*FTP, func(), error) {
				server, err := newMockServer(0)
				if err != nil {
					return nil, nil, err
				}

				host, port, _ := net.SplitHostPort(server.addr())
				portNum := 0
				fmt.Sscanf(port, "%d", &portNum)

				client, err := NewFTP(FTPConfig{
					Host:        host,
					Port:        uint(portNum),
					Username:    "user",
					Password:    "pass",
					DialTimeout: time.Second,
				})

				cleanup := func() {
					server.close()
				}

				return client, cleanup, err
			},
			wantErr:     true,
			errContains: "EOF",
		},
		{
			name: "server not responding",
			setup: func() (*FTP, func(), error) {
				server, err := newMockServer(2 * time.Second)
				if err != nil {
					return nil, nil, err
				}

				host, port, _ := net.SplitHostPort(server.addr())
				portNum := 0
				fmt.Sscanf(port, "%d", &portNum)

				client, err := NewFTP(FTPConfig{
					Host:        host,
					Port:        uint(portNum),
					Username:    "user",
					Password:    "pass",
					DialTimeout: time.Second,
				})

				cleanup := func() {
					server.close()
				}

				return client, cleanup, err
			},
			wantErr:     true,
			errContains: "EOF",
		},
		{
			name: "server unreachable",
			setup: func() (*FTP, func(), error) {
				client, err := NewFTP(FTPConfig{
					Host:        "127.0.0.1",
					Port:        65535, // Недоступный порт
					Username:    "user",
					Password:    "pass",
					DialTimeout: time.Second,
				})
				return client, func() {}, err
			},
			wantErr:     true,
			errContains: "refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup, err := tt.setup()
			require.NoError(t, err)
			defer cleanup()

			err = client.Connect(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
		})
	}
}
