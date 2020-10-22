package utils

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	SetTestStorage()
	viper.Set("file.enableLocks", true)
	viper.Set("file.lockTimeoutMs", 100)

	os.Exit(m.Run())
}
