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
	var procfile, envfile, envDir, cmdKey, etcdUrl string
	var restart bool

	//need this for etcd mode
	flag.StringVar(&etcdUrl, "etcd", "http://localhost:4001", "url of etcd")

	//neuron flags
	flag.StringVar(&envDir, "env", "default", "name of env dir")
	flag.StringVar(&cmdKey, "cmd", "", "name of cmd key")
	flag.BoolVar(&restart, "r", false, "restart instead of crashing")

	//import flags
	flag.StringVar(&procfile, "p", "Procfile", "procfile location for import")
	flag.StringVar(&envfile, "e", ".env", ".env location for import")

	flag.Parse()

	n.Etcd = etcd.NewClient([]string{etcdUrl})

	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "bootstrap":
			Bootstrap(n.Etcd, flag.Arg(1))
		case "import":
			Import(n.Etcd, procfile, envfile)
		}
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
	n.ProcId = uuid.New()
	n.Ttl = 5
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
