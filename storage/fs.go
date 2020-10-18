package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/juju/fslock"
	"github.com/spf13/viper"
	"github.com/vtex/hyper-cas/sitebuilder"
	"github.com/vtex/hyper-cas/utils"
)

type FSStorage struct {
	rootPath    string
	sitesPath   string
	siteBuilder sitebuilder.SiteBuilder
}

func NewFSStorage(siteBuilder sitebuilder.SiteBuilder) (*FSStorage, error) {
	viper.SetDefault("storage.rootPath", "/tmp/hyper-cas/storage")
	viper.SetDefault("storage.sitesPath", "/tmp/hyper-cas/sites")
	rootPath := viper.GetString("storage.rootPath")
	sitesPath := viper.GetString("storage.sitesPath")

	err := os.MkdirAll(rootPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(sitesPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &FSStorage{
		rootPath:    rootPath,
		sitesPath:   sitesPath,
		siteBuilder: siteBuilder,
	}, nil
}

func symlink(filePath, symlinkPath string) error {
	fileRelPath, err := filepath.Rel(path.Dir(symlinkPath), filePath)
	if err != nil {
		return fmt.Errorf("Failed to find relative path between %s and %s: %v", path.Dir(filePath), symlinkPath, err)
	}

	cmd := exec.Command("ln", "-sf", fileRelPath, path.Base(symlinkPath))
	cmd.Dir = path.Dir(symlinkPath)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (st *FSStorage) filePath(hash string) string {
	fileDir := path.Join(st.rootPath, "files", hash[0:2], hash[2:4])
	return path.Join(fileDir, hash)
}

func (st *FSStorage) Store(hash string, value []byte) error {
	fileDir := path.Join(st.rootPath, "files", hash[0:2], hash[2:4])
	filePath := path.Join(fileDir, hash)

	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return err
	}

	lock := fslock.New(filePath)
	err = lock.LockWithTimeout(time.Millisecond * 100)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	err = ioutil.WriteFile(filePath, []byte(value), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (st *FSStorage) Get(hash string) ([]byte, error) {
	filePath := path.Join(st.rootPath, "files", hash[0:2], hash[2:4], hash)
	if !utils.FileExists(filePath) {
		return nil, fmt.Errorf("File %s was not found!", filePath)
	}
	lock := fslock.New(filePath)
	err := lock.LockWithTimeout(time.Millisecond * 100)
	if err != nil {
		return nil, err
	}
	defer lock.Unlock()

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return dat, nil
}

func (st *FSStorage) Has(hash string) bool {
	filePath := path.Join(st.rootPath, "files", hash[0:2], hash[2:4], hash)
	return utils.FileExists(filePath)
}

func splitFile(hash string) (string, string) {
	v := strings.Split(hash, ":")
	return v[0], v[1]
}

func (st *FSStorage) StoreDistro(root string, hashes []string) error {
	dir := path.Join(st.sitesPath, fmt.Sprintf("temp%s", root))
	defer func() {
		if utils.DirExists(dir) {
			os.RemoveAll(dir)
		}
	}()
	err := st.storeDistroLinks(dir, root, hashes)
	if err != nil {
		return err
	}
	err = st.storeDistroFile(root)
	if err != nil {
		return err
	}

	finalPath := path.Join(st.sitesPath, root)
	err = os.Rename(dir, finalPath)
	if err != nil {
		return err
	}

	return nil
}

func (st *FSStorage) storeDistroLinks(dir, root string, hashes []string) error {
	for _, item := range hashes {
		filename, hash := splitFile(item)
		filePath := st.filePath(hash)
		symlinkPath := path.Join(dir, filename)
		symlinkDir := path.Dir(symlinkPath)
		err := os.MkdirAll(symlinkDir, os.ModePerm)
		if err != nil {
			return err
		}
		err = symlink(filePath, symlinkPath)
		if err != nil {
			return fmt.Errorf("Error creating symlink between %s and %s: %v", filePath, symlinkPath, err)
		}
	}

	return nil
}

func (st *FSStorage) storeDistroFile(root string) error {
	filePath := path.Join(st.rootPath, "distros", root)
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	lock := fslock.New(filePath)
	err = lock.LockWithTimeout(time.Millisecond * 100)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	err = ioutil.WriteFile(filePath, []byte(""), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (st *FSStorage) GetDistro(root string) ([]string, error) {
	filePath := path.Join(st.rootPath, "distros", root)
	lock := fslock.New(filePath)
	err := lock.LockWithTimeout(time.Millisecond * 100)
	if err != nil {
		return nil, err
	}
	defer lock.Unlock()

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var contents []string
	err = json.Unmarshal(dat, &contents)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (st *FSStorage) HasDistro(root string) bool {
	filePath := path.Join(st.rootPath, "distros", root)
	return utils.FileExists(filePath)
}

func (st *FSStorage) StoreLabel(label, hash string) error {
	err := st.storeLabelFile(label, hash)
	if err != nil {
		return err
	}
	err = st.storeLabelConf(label, hash)
	if err != nil {
		return err
	}

	return nil
}

func (st *FSStorage) storeLabelFile(label, hash string) error {
	filePath := path.Join(st.rootPath, "labels", label)
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	lock := fslock.New(filePath)
	err = lock.LockWithTimeout(time.Millisecond * 100)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return ioutil.WriteFile(filePath, []byte(hash), 0644)
}

func (st *FSStorage) storeLabelConf(label, hash string) error {
	conf, err := st.siteBuilder.Generate(label, hash)
	if err != nil {
		return err
	}

	confPath := path.Join(st.sitesPath, fmt.Sprintf("%s.conf", label))
	lock := fslock.New(confPath)
	err = lock.LockWithTimeout(time.Millisecond * 100)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return ioutil.WriteFile(confPath, []byte(conf), 0644)
}

func (st *FSStorage) GetLabel(label string) (string, error) {
	filePath := path.Join(st.rootPath, "labels", label)
	lock := fslock.New(filePath)
	err := lock.LockWithTimeout(time.Millisecond * 100)
	if err != nil {
		return "", err
	}
	defer lock.Unlock()

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

func (st *FSStorage) HasLabel(label string) bool {
	filePath := path.Join(st.rootPath, "labels", label)
	return utils.FileExists(filePath)
}
