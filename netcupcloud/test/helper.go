package test

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

// In order to test the units, you need:
// 1. Create an netcup account https://www.netcup.eu
// 2. have access to SCP, CPP and SCP Webservices (SOAP API)
// 3. Be able to pay for the services
// some tests cost money to perform,
// to enable/disable the tests use:
// enablePaidTest
const (
	local_env_path = "/test/local.env"
	projectDirName = "netcupcloud"
	enablePaidTest = false
)

type Environment struct {
	CustNo         string
	CustPwd        string
	SCPNo          string
	SCPPwd         string
	EnablePaidTest bool
}

var env Environment

// dynamicly loads the test env file
func LoadLocalEnvironment() Environment {
	if env == (Environment{}) {
		projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
		currentWorkDirectory, _ := os.Getwd()
		rootPath := projectName.Find([]byte(currentWorkDirectory))
		err := godotenv.Load(string(rootPath) + local_env_path)

		if err != nil {
			log.Fatalf("Cant load the Local Environment. Err: %s", err)
		}
		env = Environment{os.Getenv("CUST_NO"), os.Getenv("CUST_PWD"), os.Getenv("SCP_NO"), os.Getenv("SCP_PWD"), enablePaidTest}
	}
	return env
}
