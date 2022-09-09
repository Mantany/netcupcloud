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

func TestEUShop_AuthWrongPasswordWrongUsername(t *testing.T) {
	euShopSession := NewEUShop("Wrong", "password")
	err := euShopSession.auth()
	if err == nil {
		t.Error("Expected wrong password error!")
	}
}

func TestEUShop_AuthRightPasswordRightUsername(t *testing.T) {
	euShopSession := NewEUShop(env.CustNo, env.CustPwd)
	err := euShopSession.auth()
	if err != nil {
		t.Error("Expected successful authentication")
	}
}

func TestEUShop_PutIntoChartWrongId(t *testing.T) {
	euShopSession := NewEUShop(env.CustNo, env.CustPwd)
	err := euShopSession.PutIntoChart(9090989345)
	if err == nil {
		t.Error("Expected wrong ID error!")
	}
}

func TestEUShop_PutIntoChartRightId(t *testing.T) {
	euShopSession := NewEUShop(env.CustNo, env.CustPwd)
	err := euShopSession.PutIntoChart(2948)
	if err != nil {
		t.Error("Expected successful product chart")
	}
}

func TestEUShop_ReleaseOrderWithoutItemsInChart(t *testing.T) {
	euShopSession := NewEUShop(env.CustNo, env.CustPwd)
	err := euShopSession.ReleaseOrder()
	if err == nil {
		t.Error("Expected error, releasing Items without chart")
	}
}

func TestEUShop_ReleaseOrderWithItemInChart(t *testing.T) {
	euShopSession := NewEUShop(env.CustNo, env.CustPwd)
	err := euShopSession.PutIntoChart(2948)
	if err != nil {
		t.Error("Expected successful release")
	}
	erro := euShopSession.ReleaseOrder()
	if erro != nil {
		t.Error("Expected successful release")
	}
}
