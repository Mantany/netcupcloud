package session

import (
	"os"
	"testing"

	"github.com/mantany/netcupcloud/netcupcloud/test"
)

var env test.Environment

// prepare the local test environment:
func TestMain(m *testing.M) {
	env = test.LoadLocalEnvironment()
	os.Exit(m.Run())
}
