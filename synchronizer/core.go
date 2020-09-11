package synchronizer

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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

func (s *Sync) UploadFile(path, content string) (string, error) {
	hashBytes := sha256.Sum256([]byte(content))
	hash := fmt.Sprintf("%x", hashBytes)
	status, body := doReq("HEAD", fmt.Sprintf("http://localhost:2485/file/%s", hash), "")
	if status == 200 {
		fmt.Printf("* %s - Already up-to-date.\n", path)
		return hash, nil
	}
	status, body = doReq("PUT", "http://localhost:2485/file", content)
	if status != 200 {
		return "", fmt.Errorf("Failed to put %s. Status: %d Error: %s", path, status, body)
	}
	fmt.Printf("* %s - Updated (hash: %s).\n", path, body)
	return body, nil
}

func (s *Sync) UploadDistro(hashes map[string]string) (string, error) {
	keys := make([]string, 0, len(hashes))
	for key := range hashes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, path := range keys {
		sb.WriteString(path)
		sb.WriteString(":")
		sb.WriteString(hashes[path])
		sb.WriteString("\n")
	}
	content := sb.String()
	status, body := doReq("PUT", "http://localhost:2485/distro", content)
	if status != 200 {
		return "", fmt.Errorf("Failed to put new distro. Status: %d Error: %s", status, body)
	}
	fmt.Printf("* Distro %s is up-to-date.\n", body)
	return body, nil
}

func (s *Sync) Run() (string, error) {
	files, contents, err := listFiles(s.rootDir)
	if err != nil {
		panic(err)
	}
	hashes := map[string]string{}
	for i, path := range files {
		content := contents[i]
		hash, err := s.UploadFile(path, content)
		if err != nil {
			return "", err
		}
		hashes[path] = hash
	}
	distro, err := s.UploadDistro(hashes)
	if err != nil {
		return "", err
	}
	return distro, nil
}
