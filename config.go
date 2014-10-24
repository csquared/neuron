package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/coreos/go-etcd/etcd"
)

func Config(n *Neuron) {
	var envDir, cmdKey, etcdUrl string

	flag.StringVar(&envDir, "env", "default", "name of env dir")
	flag.StringVar(&cmdKey, "cmd", "", "name of cmd key")
	flag.StringVar(&etcdUrl, "etcd", "http://localhost:4001", "url of etcd")

	flag.Parse()

	if envDir == "" || cmdKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	n.Etcd = etcd.NewClient([]string{etcdUrl})
	envDir, cmdKey = resolveKeys(envDir, cmdKey)
	n.EnvDir = envDir
	n.Env = GetEnv(n.Etcd, envDir)
	n.CmdKey = cmdKey
	n.Command = GetCmd(n.Etcd, cmdKey)
	n.AppName = appName()
}

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
