package forwardproxy

import (
	"context"
	"fmt"
	"github.com/achillescres/pkg/utils"
	"github.com/elazarl/goproxy"
	"net"
	"net/http"
	"net/url"
)

type ForwardProxy struct {
	port      int
	server    *goproxy.ProxyHttpServer
	forwardTo url.URL
}

func New(port int, forwardTo *url.URL) *ForwardProxy {
	server := goproxy.NewProxyHttpServer()
	//server.Verbose = true

	fp := &ForwardProxy{
		server:    server,
		forwardTo: *forwardTo,
		port:      port,
	}
	server.Tr.Proxy = func(*http.Request) (*url.URL, error) {
		return &fp.forwardTo, nil
	}

	return fp
}

func (fp *ForwardProxy) Port() int {
	return fp.port
}

func (fp *ForwardProxy) Forward(forwardTo *url.URL) {
	fp.forwardTo = *forwardTo
}

func (fp *ForwardProxy) Run(ctx context.Context) error {
	ew := utils.NewErrorWrapper("ForwardProxy - Run")

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", fp.port),
		Handler: fp.server,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	return ew(server.ListenAndServe())
}
