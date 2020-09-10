package synchronizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

type Sync struct {
	rootDir string
}

func doReq(method, url, body string) (int, string) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	req.Header.SetMethodBytes([]byte(method))
	req.SetRequestURI(url)
	if body != "" {
		req.SetBodyString(body)
	}

	fasthttp.Do(req, resp)

	status := resp.StatusCode()
	bodyBytes := resp.Body()
	return status, string(bodyBytes)
}

func readAll(path string) (string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func listFiles(path string) ([]string, []string, error) {
	files := []string{}
	contents := []string{}
	err := filepath.Walk(path,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, strings.Replace(p, path+"/", "", 1))
				content, err := readAll(p)
				if err != nil {
					return err
				}
				contents = append(contents, content)
			}
			return nil
		})
	if err != nil {
		return nil, nil, err
	}

	return files, contents, err
}

func NewSync(root string) *Sync {
	return &Sync{rootDir: root}
}

func (s *Sync) Run() error {
	files, contents, err := listFiles(s.rootDir)
	if err != nil {
		panic(err)
	}
	for i, path := range files {
		content := contents[i]
		// TODO: Implement requests
		fmt.Println(path, content)
	}
	return nil
}
