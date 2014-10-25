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
	var restart bool

	flag.StringVar(&envDir, "env", "default", "name of env dir")
	flag.StringVar(&cmdKey, "cmd", "", "name of cmd key")
	flag.StringVar(&etcdUrl, "etcd", "http://localhost:4001", "url of etcd")
	flag.BoolVar(&restart, "r", false, "restart instead of crashing")

	flag.Parse()

	n.Etcd = etcd.NewClient([]string{etcdUrl})

	if len(os.Args) == 2 && os.Args[1] == "bootstrap" {
		Bootstrap(n.Etcd, flag.Arg(1))
		os.Exit(0)
	}

	if len(os.Args) == 2 && os.Args[1] == "import" {
		Import(n.Etcd)
		os.Exit(0)
	}

	if envDir == "" || cmdKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	envDir, cmdKey = resolveKeys(envDir, cmdKey)
	n.EnvDir = envDir
	n.Env = GetEnv(n.Etcd, envDir)
	n.CmdKey = cmdKey
	n.Command = GetCmd(n.Etcd, cmdKey)
	n.AppName = appName()
	n.Restart = restart
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
