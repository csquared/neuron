package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"code.google.com/p/go-uuid/uuid"

	"github.com/coreos/go-etcd/etcd"
)

func Config(n *Neuron) {
	var appName, envName, cmdKey, etcdUrl string
	var restart bool

	//need this for etcd mode
	flag.StringVar(&etcdUrl, "etcd", "http://localhost:4001", "url of etcd")

	//neuron flags
	flag.StringVar(&appName, "app", GetAppName(), "name of the app, usually the directory")
	flag.StringVar(&envName, "env", "dev", "name of env")
	flag.StringVar(&cmdKey, "cmd", "", "name of command")
	flag.BoolVar(&restart, "r", false, "restart instead of crashing")

	flag.Parse()

	n.Etcd = etcd.NewClient([]string{etcdUrl})

	if appName == "" {
		appName = GetAppName()
	}

	n.AppName = appName
	n.ProcId = uuid.New()
	n.EnvName = envName
	n.CmdName = cmdKey
	n.Restart = restart
	n.Ttl = 5
}

func GetAppName() string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Base(workingDir)
}
