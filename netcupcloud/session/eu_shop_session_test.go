package session

import (
	"os"
	"testing"

	"github.com/Mantany/netcupcloud/netcupcloud/test"
)

var env test.Environment

// prepare the local test environment:
func TestMain(m *testing.M) {
	env = test.LoadLocalEnvironment()
	os.Exit(m.Run())
}

func TestEUShopAuthWrongPasswordWrongUsername(t *testing.T) {
	eu_shop_session := NewEUShopSession("Wrong", "password")
	err := eu_shop_session.auth()
	if err == nil {
		t.Error("Expected wrong password error!")
	}
}

func TestEUShopAuthRightPasswordRightUsername(t *testing.T) {
	eu_shop_session := NewEUShopSession(env.Cust_no, env.Cust_pwd)
	err := eu_shop_session.auth()
	if err != nil {
		t.Error("Expected successful authentication")
	}
}
