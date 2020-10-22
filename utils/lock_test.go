package utils

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func testFile() (string, func()) {
	file, err := ioutil.TempFile("", "test")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(file.Name(), []byte("test"), 0644)
	return file.Name(), func() {
		os.Remove(file.Name())
	}
}

func TestLockWhenLockEnabled(t *testing.T) {
	viper.Set("file.enableLocks", true)
	viper.Set("file.lockTimeoutMs", 20)
	f, cleanup := testFile()
	defer cleanup()

	unlock1, err := Lock(f)
	defer unlock1()

	assert.NoError(t, err)
}

func TestLockFailsWhenAlreadyLocked(t *testing.T) {
	viper.Set("file.enableLocks", true)
	viper.Set("file.lockTimeoutMs", 20)
	f, cleanup := testFile()
	defer cleanup()
	unlock1, err := Lock(f)
	defer unlock1()
	assert.NoError(t, err)

	unlock2, err := Lock(f)
	defer unlock2()

	assert.Error(t, err)
}

func TestLockIsNoOpWhenDisabled(t *testing.T) {
	viper.Set("file.enableLocks", false)
	viper.Set("file.lockTimeoutMs", 20)
	f, cleanup := testFile()
	defer cleanup()
	unlock1, err := Lock(f)
	defer unlock1()

	unlock2, err := Lock(f)
	defer unlock2()

	assert.NoError(t, err)
}
