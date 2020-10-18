package serve

import (
	"os"
	"testing"

	"github.com/vtex/hyper-cas/utils"
)

func TestMain(m *testing.M) {
	utils.SetTestStorage()
	os.Exit(m.Run())
}
