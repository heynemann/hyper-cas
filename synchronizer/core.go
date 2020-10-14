package synchronizer

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/valyala/fasthttp"
)

type Sync struct {
	rootDir string
	apiURL  string
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

func NewSync(root, apiURL string) *Sync {
	return &Sync{rootDir: root, apiURL: apiURL}
}

func (s *Sync) doReq(method, reqUrl, body string) (int, string) {
	u, err := url.Parse(s.apiURL)
	if err != nil {
		return 500, fmt.Sprintf("Invalid URL %s", s.apiURL)
	}
	u.Path = path.Join(u.Path, reqUrl)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	req.Header.SetMethodBytes([]byte(method))
	req.SetRequestURI(u.String())
	if body != "" {
		req.SetBodyString(body)
	}

	err = fasthttp.Do(req, resp)
	if err != nil {
		return 500, fmt.Sprintf("%s for %s failed with %v.", method, reqUrl, err)
	}

	status := resp.StatusCode()
	bodyBytes := resp.Body()
	return status, string(bodyBytes)
}

func (s *Sync) UploadFile(path, content string) (string, error) {
	hashBytes := sha1.Sum([]byte(content))
	hash := fmt.Sprintf("%x", hashBytes)
	fileURL := fmt.Sprintf("/file/%s", hash)
	status, body := s.doReq("HEAD", fileURL, "")
	if status == 200 {
		fmt.Printf("* %s - Already up-to-date.\n", path)
		return hash, nil
	}
	status, body = s.doReq("PUT", "/file", content)
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
	status, body := s.doReq("PUT", "/distro", content)
	if status != 200 {
		return "", fmt.Errorf("Failed to put new distro. Status: %d Error: %s", status, body)
	}
	fmt.Printf("* Distro %s is up-to-date.\n", body)
	return body, nil
}

func (s *Sync) SetLabel(label, hash string) error {
	status, body := s.doReq("PUT", "/label", fmt.Sprintf("label=%s&hash=%s", label, hash))
	if status != 200 {
		return fmt.Errorf("Failed to put new distro. Status: %d Error: %s", status, body)
	}
	fmt.Printf("* Updated label %s => %s.\n", label, hash)
	return nil

}

func (s *Sync) Run(label string) (string, error) {
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
	if label != "" {
		err = s.SetLabel(label, distro)
		if err != nil {
			return "", err
		}
	}
	return distro, nil
}
