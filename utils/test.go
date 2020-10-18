package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

type App interface {
	GetRouter() *router.Router
}

// serve serves http request using provided fasthttp handler
func serveRequest(app App, req *http.Request) (*http.Response, error) {
	router := app.GetRouter()
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, router.Handler)
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

func DoRequest(app App, method, url, body string) (*http.Response, int, string, error) {
	var bodyReader io.Reader
	if method != "GET" && body != "" {
		bodyReader = bytes.NewBuffer([]byte(body))
	}
	r, err := http.NewRequest(method, fmt.Sprintf("http://localhost/%s", url), bodyReader)
	if err != nil {
		return nil, 500, "", err
	}

	res, err := serveRequest(app, r)
	if err != nil {
		return res, 500, "", err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, 500, "", err
	}

	return res, res.StatusCode, string(resBody), nil
}
