package session

import (
	"testing"
)

func TestEUShop_AuthWrongPasswordWrongUsername(t *testing.T) {
	euShopSession := NewEUShop("Wrong", "password")
	err := euShopSession.auth()
	if err == nil {
		t.Error("Expected wrong password error!")
	}
}

func TestEUShop_AuthRightPasswordRightUsername(t *testing.T) {
	euShopSession := NewEUShop(env.CCPNo, env.CCPPwd)
	err := euShopSession.auth()
	if err != nil {
		t.Error("Expected successful authentication, please make sure that you have set the environment variables correctly")
	}
}

func TestEUShop_PutIntoChartWrongId(t *testing.T) {
	euShopSession := NewEUShop(env.CCPNo, env.CCPPwd)
	err := euShopSession.PutIntoChart(9090989345)
	if err == nil {
		t.Error("Expected wrong ID error!")
	}
}

func TestEUShop_PutIntoChartRightId(t *testing.T) {
	euShopSession := NewEUShop(env.CCPNo, env.CCPPwd)
	err := euShopSession.PutIntoChart(2948)
	if err != nil {
		t.Error("Expected successful product chart")
	}
}

func TestEUShop_ReleaseOrderWithoutItemsInChart(t *testing.T) {
	euShopSession := NewEUShop(env.CCPNo, env.CCPPwd)
	err := euShopSession.ReleaseOrder()
	if err == nil {
		t.Error("Expected error, releasing Items without chart")
	}
}

func TestEUShop_ReleaseOrderWithItemInChart(t *testing.T) {
	if env.EnablePaidTest {
		euShopSession := NewEUShop(env.CCPNo, env.CCPPwd)
		err := euShopSession.PutIntoChart(2948)
		if err != nil {
			t.Error("Expected successful release")
		}
		erro := euShopSession.ReleaseOrder()
		if erro != nil {
			t.Error("Expected successful release")
		}
	}
}
