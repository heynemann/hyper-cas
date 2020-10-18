package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

type App interface {
	HandleError(func(ctx *fasthttp.RequestCtx) error) func(ctx *fasthttp.RequestCtx)
}

// serve serves http request using provided fasthttp handler
func serveRequest(app App, handler func(ctx *fasthttp.RequestCtx) error, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, app.HandleError(handler))
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}

func DoRequest(app App, handler func(ctx *fasthttp.RequestCtx) error, method, url, body string) (*http.Response, int, string, error) {
	var bodyReader io.Reader
	if method != "GET" && body != "" {
		bodyReader = bytes.NewBuffer([]byte(body))
	}
	r, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 500, "", err
	}

	res, err := serveRequest(app, handler, r)
	if err != nil {
		return res, 500, "", err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, 500, "", err
	}

	return res, res.StatusCode, string(resBody), nil
}
