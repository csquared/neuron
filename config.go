package main

import (
	"log"
	"os"
	"path/filepath"
)

func resolveKeys(envDir, cmdKey string) (string, string) {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if envDir == "" || cmdKey == "" {
		log.Fatal("you need to supply arguments")
	}

	appName := filepath.Base(workingDir)

	if envDir[0] != '/' {
		envDir = "/services/" + appName + "/envs/" + envDir
	}
	if cmdKey[0] != '/' {
		cmdKey = "/services/" + appName + "/processes/" + cmdKey
	}

	return envDir, cmdKey
}
