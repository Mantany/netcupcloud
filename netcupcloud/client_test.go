package netcupcloud

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

const (
	local_env_name = "local.env"
)

type Environment struct {
	cust_no  string
	cust_pwd string
}

func loadLocalEnvironment() *Environment {
	err := godotenv.Load(local_env_name)
	if err != nil {
		log.Fatalf("Cant load the Local Environment. Err: %s", err)
	}
	env := &Environment{os.Getenv("CUST_NO"), os.Getenv("CUST_PWD")}
	return env
}

func TestEUShopAuthWrongPasswordWrongUsername(t *testing.T) {
	test_client := NewClient("Wrong", "password")
	err := test_client.eu_shop_session.Auth()
	if err == nil {
		t.Error("Expected wrong password error! ")
	}
}

func TestEUShopAuthRightPasswordRightUsername(t *testing.T) {
	env := loadLocalEnvironment()
	test_client := NewClient(env.cust_no, env.cust_pwd)
	err := test_client.eu_shop_session.Auth()
	if err != nil {
		t.Error("Expected finished auth")
	}
}
