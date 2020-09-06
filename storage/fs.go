package storage

import (
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

func (st *FSStorage) Store(hash string, value []byte) error {
	fileDir := path.Join(st.rootPath, hash[0:2], hash[2:4])
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
	filePath := path.Join(st.rootPath, hash[0:2], hash[2:4], hash)
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
