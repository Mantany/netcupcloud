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

func TestEUShopAuthWrongPasswordWrongUsername(t *testing.T) {
	euShopSession := NewEUShopSession("Wrong", "password")
	err := euShopSession.auth()
	if err == nil {
		t.Error("Expected wrong password error!")
	}
}

func TestEUShopAuthRightPasswordRightUsername(t *testing.T) {
	euShopSession := NewEUShopSession(env.CustNo, env.CustPwd)
	err := euShopSession.auth()
	if err != nil {
		t.Error("Expected successful authentication")
	}
}
