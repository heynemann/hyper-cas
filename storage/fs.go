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
)

type FSStorage struct {
	rootPath    string
	sitesPath   string
	siteBuilder sitebuilder.SiteBuilder
}

func NewFSStorage(siteBuilder sitebuilder.SiteBuilder) (*FSStorage, error) {
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func symlink(filePath, symlinkPath string) error {
	// symlinkPathTmp := symlinkPath + ".tmp"
	// if err := os.Remove(symlinkPathTmp); err != nil && !os.IsNotExist(err) {
	// return err
	// }

	// if err := os.Symlink(filePath, symlinkPathTmp); err != nil {
	// return err
	// }
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

	// if err := os.Rename(symlinkPathTmp, symlinkPath); err != nil {
	// return err
	// }
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
	return fileExists(filePath)
}

func splitFile(hash string) (string, string) {
	v := strings.Split(hash, ":")
	return v[0], v[1]
}

func (st *FSStorage) StoreDistro(root string, hashes []string) error {
	for _, item := range hashes {
		filename, hash := splitFile(item)
		filePath := st.filePath(hash)
		symlinkPath := path.Join(st.sitesPath, root, filename)
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

func (st *FSStorage) GetDistro(root string) ([]string, error) {
	filePath := path.Join(st.rootPath, "distros", root[0:2], root[2:4], root)
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

func (st *FSStorage) HasDistro(hash string) bool {
	filePath := path.Join(st.rootPath, "distros", hash[0:2], hash[2:4], hash)
	return fileExists(filePath)
}

func (st *FSStorage) StoreLabel(label, hash string) error {
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

	err = ioutil.WriteFile(filePath, []byte(hash), 0644)
	if err != nil {
		return err
	}
	conf, err := st.siteBuilder.Generate(label, hash)
	if err != nil {
		return err
	}

	confPath := path.Join(st.sitesPath, fmt.Sprintf("%s.conf", label))
	lock = fslock.New(confPath)
	err = lock.LockWithTimeout(time.Millisecond * 100)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	err = ioutil.WriteFile(confPath, []byte(conf), 0644)
	if err != nil {
		return err
	}

	return nil
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
	return fileExists(filePath)
}
