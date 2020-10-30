package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
	"github.com/vtex/hyper-cas/sitebuilder"
	"github.com/vtex/hyper-cas/utils"
	"go.uber.org/zap"
)

// FSStorage for keeping all the CAS data in the filesystem
type FSStorage struct {
	rootPath    string
	sitesPath   string
	siteBuilder sitebuilder.SiteBuilder
}

// NewFSStorage with the specified settings
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
	symlinkPathTmp := symlinkPath + ".tmp"
	logger := utils.LoggerWith(
		zap.String("filePath", filePath),
		zap.String("symlinkPath", symlinkPath),
	)
	if err := os.Remove(symlinkPathTmp); err != nil && !os.IsNotExist(err) {
		logger.Error("failed to remove previous symlink", zap.Error(err))
		return err
	}

	if err := os.Symlink(filePath, symlinkPathTmp); err != nil {
		logger.Error(
			"failed to create temporary symlink",
			zap.String("tempSymlink", symlinkPathTmp),
			zap.Error(err),
		)
		return err
	}

	if err := os.Rename(symlinkPathTmp, symlinkPath); err != nil {
		logger.Error(
			"failed to move temporary symlink to symlink",
			zap.String("tempSymlink", symlinkPathTmp),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (st *FSStorage) filePath(hash string) string {
	fileDir := path.Join(st.rootPath, "files", hash[0:2], hash[2:4])
	return path.Join(fileDir, hash)
}

// Store files in the filesystem
func (st *FSStorage) Store(hash string, value []byte) error {
	fileDir := path.Join(st.rootPath, "files", hash[0:2], hash[2:4])
	filePath := path.Join(fileDir, hash)

	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return err
	}

	unlock, err := utils.Lock(filePath)
	defer unlock()

	err = ioutil.WriteFile(filePath, []byte(value), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Get a file from the filesystem
func (st *FSStorage) Get(hash string) ([]byte, error) {
	filePath := path.Join(st.rootPath, "files", hash[0:2], hash[2:4], hash)
	if !utils.FileExists(filePath) {
		return nil, fmt.Errorf("file %s was not found", filePath)
	}
	unlock, err := utils.Lock(filePath)
	defer unlock()

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return dat, nil
}

// Has the file in the filesystem?
func (st *FSStorage) Has(hash string) bool {
	filePath := path.Join(st.rootPath, "files", hash[0:2], hash[2:4], hash)
	return utils.FileExists(filePath)
}

func splitFile(hash string) (string, string) {
	v := strings.Split(hash, ":")
	return v[0], v[1]
}

// StoreDistro in the filesytem
func (st *FSStorage) StoreDistro(root string, hashes []string) error {
	dir := path.Join(st.sitesPath, fmt.Sprintf("%s%s", utils.RandString(32), root))
	defer func() {
		if utils.DirExists(dir) {
			os.RemoveAll(dir)
		}
	}()
	err := st.storeDistroLinks(dir, root, hashes)
	if err != nil {
		return err
	}
	err = st.storeDistroFile(root, hashes)
	if err != nil {
		return err
	}

	finalPath := path.Join(st.sitesPath, root)
	err = os.Rename(dir, finalPath)
	if err != nil {
		sErr, ok := err.(*os.LinkError)
		if ok && sErr.Op == "rename" {
			utils.LogWarn("Distribution path already exist. Ignoring rename...", zap.Error(err), zap.String("path", finalPath))
			return nil
		}
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

func (st *FSStorage) storeDistroFile(root string, hashes []string) error {
	filePath := path.Join(st.rootPath, "distros", root)
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	unlock, err := utils.Lock(filePath)
	defer unlock()

	contents, err := json.Marshal(hashes)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, contents, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetDistro from the filesystem
func (st *FSStorage) GetDistro(root string) ([]string, error) {
	filePath := path.Join(st.rootPath, "distros", root)
	unlock, err := utils.Lock(filePath)
	defer unlock()

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

// HasDistro in the filesystem?
func (st *FSStorage) HasDistro(root string) bool {
	filePath := path.Join(st.rootPath, "distros", root)
	return utils.FileExists(filePath)
}

// StoreLabel in the filesystem
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

	unlock, err := utils.Lock(filePath)
	defer unlock()

	return ioutil.WriteFile(filePath, []byte(hash), 0644)
}

func (st *FSStorage) storeLabelConf(label, hash string) error {
	conf, err := st.siteBuilder.Generate(label, hash)
	if err != nil {
		return err
	}

	confPath := path.Join(st.sitesPath, fmt.Sprintf("%s.conf", label))
	unlock, err := utils.Lock(confPath)
	defer unlock()

	return ioutil.WriteFile(confPath, []byte(conf), 0644)
}

// GetLabel from the filesystem
func (st *FSStorage) GetLabel(label string) (string, error) {
	filePath := path.Join(st.rootPath, "labels", label)
	unlock, err := utils.Lock(filePath)
	defer unlock()

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

// HasLabel in the filesystem?
func (st *FSStorage) HasLabel(label string) bool {
	filePath := path.Join(st.rootPath, "labels", label)
	return utils.FileExists(filePath)
}
