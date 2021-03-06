package utils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/fasthttp/router"
	"github.com/spf13/viper"
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
		bodyReader = strings.NewReader(body)
	}
	r, err := http.NewRequest(method, fmt.Sprintf("http://localhost/%s", url), bodyReader)
	if err != nil {
		return nil, 500, "", err
	}
	if method == "POST" || method == "PUT" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

func SetTestStorage() {
	os.RemoveAll("/tmp/hyper-cas-test")
	viper.Set("storage.rootPath", "/tmp/hyper-cas-test/storage")
	viper.Set("storage.sitesPath", "/tmp/hyper-cas-test/sites")
}
