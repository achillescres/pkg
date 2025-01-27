package forwardproxyfabric

import (
	"context"
	"fmt"
	"github.com/elazarl/goproxy"
	"net"
	"net/http"
	"net/url"
)

type Server struct {
	close   func() error
	forward url.URL
	port    int
}

func (s Server) Reload() {

}

func (s Server) Close() error {
	return s.close()
}

func (s Server) Forward(to url.URL) {
	s.forward = to
}

func (s Server) Port() int {
	return s.port
}

func (s Server) Url() *url.URL {
	u, err := url.Parse(fmt.Sprintf("http://localhost:%d", s.Port()))
	if err != nil {
		panic(err)
	}
	return u
}

func New(ctx context.Context, serverErrCallback func(error), forward url.URL) *Server {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic("listen tcp :0: " + err.Error())
	}

	sc := Server{
		close:   nil,
		forward: forward,
		port:    listener.Addr().(*net.TCPAddr).Port,
	}
	fpServer := goproxy.NewProxyHttpServer()
	fpServer.Tr.Proxy = func(*http.Request) (*url.URL, error) { return &sc.forward, nil }
	server := &http.Server{
		Handler: fpServer,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	sc.close = server.Close
	go func() {
		if err := server.Serve(listener); err != nil {
			serverErrCallback(err)
		}
	}()

	return &sc
}

func NewWithPort(ctx context.Context, port int, forward url.URL, serverErrCallback func(error)) *Server {
	sc := Server{
		close:   nil,
		forward: forward,
		port:    port,
	}
	fpServer := goproxy.NewProxyHttpServer()
	fpServer.Tr.Proxy = func(*http.Request) (*url.URL, error) { return &sc.forward, nil }
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: fpServer,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	sc.close = server.Close
	go func() {
		if err := server.ListenAndServe(); err != nil {
			serverErrCallback(err)
		}
	}()

	return &sc
}
