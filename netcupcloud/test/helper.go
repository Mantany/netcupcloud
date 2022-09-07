package test

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const (
	local_env_path = "/test/local.env"
	projectDirName = "netcupcloud"
)

type Environment struct {
	Cust_no  string
	Cust_pwd string
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
		env = Environment{os.Getenv("CUST_NO"), os.Getenv("CUST_PWD")}
	}
	return env
}
