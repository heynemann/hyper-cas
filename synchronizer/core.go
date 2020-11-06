package synchronizer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
	"github.com/vtex/hyper-cas/utils"
	"go.uber.org/zap"
)

// Sync contents of root dir to CAS
type Sync struct {
	rootDir               string
	apiURL                string
	wg                    sync.WaitGroup
	jobChan               chan *fileUpdateJob
	respChan              chan *fileUpdateResponse
	fileUploadClient      *httpclient.Client
	metadataClient        *httpclient.Client
	maxConcurrentRequests int
}

// NewSync creates a Sync
func NewSync(root, apiURL string, requestRetriesCount, maxConcurrentRequests, httpTimeoutMs, distroHTTPTimeoutMs int) *Sync {
	fileUploadClient := initHTTPClient(requestRetriesCount, maxConcurrentRequests, httpTimeoutMs)
	metadataClient := initHTTPClient(requestRetriesCount, maxConcurrentRequests, distroHTTPTimeoutMs)
	s := &Sync{
		rootDir:               root,
		apiURL:                apiURL,
		fileUploadClient:      fileUploadClient,
		metadataClient:        metadataClient,
		maxConcurrentRequests: maxConcurrentRequests,
	}
	return s
}

type fileUpdateJob struct {
	path     string
	filePath string
}

type fileUpdateResponse struct {
	path          string
	hash          string
	duration      time.Duration
	alreadyExists bool
}

func readAll(path string) (string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func listFiles(jobChan chan *fileUpdateJob, path string) (int, error) {
	files := 0
	err := filepath.Walk(
		path,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files++
				utils.LogDebug("added job to queue", zap.String("path", p))
				jobChan <- &fileUpdateJob{
					path: p,
				}
			}
			return nil
		},
	)
	return files, err
}

func initHTTPClient(requestRetryCount, maxConcurrentRequests, httpTimeoutMs int) *httpclient.Client {
	// First set a backoff mechanism. Constant backoff increases the backoff at a constant rate
	backoffInterval := 2 * time.Millisecond
	// Define a maximum jitter interval. It must be more than 1*time.Millisecond
	maximumJitterInterval := 5 * time.Millisecond

	backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	return httpclient.NewClient(
		httpclient.WithHTTPTimeout(time.Duration(httpTimeoutMs)*time.Millisecond),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(requestRetryCount),
	)
}

func (s *Sync) startWorkers(workerCount int) {
	s.jobChan = make(chan *fileUpdateJob, workerCount)
	s.respChan = make(chan *fileUpdateResponse, 1000000000)

	for i := 0; i < workerCount; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

func (s *Sync) worker() {
	defer func() {
		utils.LogDebug("worker closed")
		(&s.wg).Done()
	}()

	for {
		job, ok := <-s.jobChan
		if !ok {
			return
		}
		utils.LogDebug("processing job", zap.String("filePath", job.filePath), zap.String("path", job.path))
		filePath := strings.Replace(job.path, s.rootDir+"/", "", 1)
		content, err := readAll(job.path)
		logger := utils.LoggerWith(
			zap.String("path", filePath),
			zap.String("filePath", job.path),
		)
		if err != nil {
			logger.Error("failed to read file.", zap.Error(err))
			s.respChan <- nil
			continue
		}
		hash, alreadyExists, duration, err := s.uploadFile(filePath, content)
		if err != nil {
			logger.Error("failed to upload file.", zap.String("path", job.path), zap.Error(err))
			s.respChan <- nil
			continue
		}
		s.respChan <- &fileUpdateResponse{
			path:          filePath,
			hash:          hash,
			duration:      duration,
			alreadyExists: alreadyExists,
		}
	}
}

func (s *Sync) doReq(client *httpclient.Client, method, reqURL, body string, isURLEncoded bool) (int, string) {
	u, err := url.Parse(s.apiURL)
	if err != nil {
		return 500, fmt.Sprintf("Invalid URL %s", s.apiURL)
	}
	u.Path = path.Join(u.Path, reqURL)

	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(method, u.String(), bytes.NewBuffer([]byte(body)))
	} else {
		req, _ = http.NewRequest(method, u.String(), nil)
	}
	if isURLEncoded {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := client.Do(req)
	if err != nil {
		return 500, fmt.Sprintf("%s for %s failed with %v.", method, reqURL, err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 500, fmt.Sprintf("%s for %s failed with %v.", method, reqURL, err)
	}

	return resp.StatusCode, string(respBody)
}

func (s *Sync) uploadFile(path, content string) (string, bool, time.Duration, error) {
	hashBytes := utils.Hash(content)
	hash := fmt.Sprintf("%x", hashBytes)
	fileURL := fmt.Sprintf("/file/%s", hash)
	start := time.Now()
	status, _ := s.doReq(s.fileUploadClient, "HEAD", fileURL, "", false)
	if status == 200 {
		return hash, true, time.Since(start), nil
	}
	status, body := s.doReq(s.fileUploadClient, "PUT", "/file", content, false)
	if status != 200 {
		return "", false, time.Since(start), fmt.Errorf("failed to put %s. Status: %d Error: %s", path, status, body)
	}
	return body, false, time.Since(start), nil
}

func (s *Sync) uploadDistro(hashes map[string]string) (string, time.Duration, error) {
	start := time.Now()
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
	status, body := s.doReq(s.metadataClient, "PUT", "/distro", content, false)
	if status != 200 {
		return "", time.Since(start), fmt.Errorf("failed to put new distro. Status: %d Error: %s", status, body)
	}
	return body, time.Since(start), nil
}

// HasDistro in hyper-cas with specified hash?
func (s *Sync) HasDistro(hash string) bool {
	status, _ := s.doReq(s.metadataClient, "HEAD", fmt.Sprintf("/distro/%s", hash), "", false)
	return status == 200
}

// SetLabel to specified hash
func (s *Sync) SetLabel(label, hash string) error {
	status, body := s.doReq(s.metadataClient, "PUT", "/label", fmt.Sprintf("label=%s&hash=%s", label, hash), true)
	if status != 200 {
		return fmt.Errorf("failed to put new distro. Status: %d Error: %s", status, body)
	}
	return nil
}

// Run the sync
func (s *Sync) Run(label string) (map[string]interface{}, error) {
	start := time.Now()
	result := map[string]interface{}{
		"timestamp": int32(time.Now().Unix()),
		"files":     []map[string]interface{}{},
		"distro":    map[string]interface{}{},
		"label": map[string]interface{}{
			"label": "",
			"hash":  "",
		},
	}
	s.startWorkers(s.maxConcurrentRequests)
	utils.LogDebug("workers started.", zap.Int("workerCount", s.maxConcurrentRequests))
	fileCount, err := listFiles(s.jobChan, s.rootDir)
	if err != nil {
		panic(err)
	}
	close(s.jobChan)
	defer close(s.respChan)
	utils.LogDebug("Waiting for jobs to finish...")
	s.wg.Wait()
	utils.LogDebug("Workers finished successfully. Reading results from response channel...")
	hashes := map[string]string{}
	for a := 0; a < fileCount; a++ {
		res := <-s.respChan
		if res == nil {
			continue
		}
		hashes[res.path] = res.hash
		result["files"] = append(result["files"].([]map[string]interface{}), map[string]interface{}{
			"path":     res.path,
			"hash":     res.hash,
			"exists":   res.alreadyExists,
			"duration": res.duration.Milliseconds(),
		})
	}
	utils.LogDebug("Hashes calculated.", zap.Int("hashes", len(hashes)))
	if len(hashes) != fileCount {
		utils.LogError(
			"failed to upload files to hyper-cas.",
			zap.Int("uploadedFiled", len(hashes)),
			zap.Int("filesToUpload", fileCount),
		)
		return nil, fmt.Errorf("failed to upload files to hyper-cas")
	}

	distro, distroDuration, err := s.uploadDistro(hashes)
	if err != nil {
		utils.LogError("distro could not be updated.", zap.Error(err))
		return nil, err
	}
	result["distro"] = map[string]interface{}{
		"hash":     distro,
		"duration": distroDuration.Milliseconds(),
	}
	utils.LogDebug("Distro updated successfully.", zap.String("distro", distro))
	if label != "" {
		utils.LogDebug("Label should be set.", zap.String("label", label), zap.String("distro", distro))
		err = s.SetLabel(label, distro)
		if err != nil {
			utils.LogError("failed to update label.", zap.String("label", label), zap.String("distro", distro), zap.Error(err))
			return nil, err
		}
		result["label"] = map[string]interface{}{
			"label": label,
			"hash":  distro,
		}
		utils.LogDebug("Label updated successfully.", zap.String("label", label), zap.String("distro", distro))
	}

	result["duration"] = time.Since(start).Milliseconds()

	return result, nil
}
