package session

import (
	"fmt"
	"os"
	"testing"

	"github.com/mantany/netcupcloud/netcupcloud/test"
)

var env test.Environment

// prepare the local test environment:
func TestMain(m *testing.M) {
	fmt.Println("Preparing the test environment with the environment vars")
	env = test.LoadLocalEnvironment()
	os.Exit(m.Run())
}
