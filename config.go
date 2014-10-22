package main

import (
	"log"
	"os"
	"path/filepath"
)

func appName() string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Base(workingDir)
}

func resolveKeys(envDir, cmdKey string) (string, string) {
	appName := appName()

	if envDir[0] != '/' {
		envDir = "/services/" + appName + "/envs/" + envDir
	}
	if cmdKey[0] != '/' {
		cmdKey = "/services/" + appName + "/processes/" + cmdKey
	}

	return envDir, cmdKey
}
