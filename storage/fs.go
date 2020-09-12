package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/juju/fslock"
	"github.com/spf13/viper"
)

type FSStorage struct {
	rootPath string
}

func NewFSStorage() (*FSStorage, error) {
	rootPath := viper.GetString("storage.rootPath")
	return &FSStorage{rootPath: rootPath}, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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

func (st *FSStorage) StoreDistro(root string, hashes []string) error {
	contents, err := json.Marshal(hashes)
	if err != nil {
		return nil
	}
	fileDir := path.Join(st.rootPath, "distros", root[0:2], root[2:4])
	filePath := path.Join(fileDir, root)

	err = os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return err
	}

	lock := fslock.New(filePath)
	err = lock.LockWithTimeout(time.Millisecond * 100)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	err = ioutil.WriteFile(filePath, contents, 0644)
	if err != nil {
		return err
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
	fileDir := path.Join(st.rootPath, "labels")
	filePath := path.Join(fileDir, label)

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

	err = ioutil.WriteFile(filePath, []byte(hash), 0644)
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
