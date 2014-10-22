package main

import (
	"flag"
	"os"

	"github.com/coreos/go-etcd/etcd"
)

func main() {
	var envDir string
	flag.StringVar(&envDir, "env", "default", "name of env dir")
	var cmdKey string
	flag.StringVar(&cmdKey, "cmd", "", "name of cmd key")
	var etcdUrl string
	flag.StringVar(&etcdUrl, "etcd", "http://localhost:4001", "url of etcd")
	flag.Parse()

	client := etcd.NewClient([]string{etcdUrl})

	if flag.Arg(0) == "bootstrap" {
		bootstrap(client, flag.Arg(1))
		return
	}

	if envDir == "" || cmdKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	envDir, cmdKey = resolveKeys(envDir, cmdKey)

	env := getEnv(client, envDir)
	command := getCmd(client, cmdKey)
	cmd := spawnProc(env, command)
	for watch(client, envDir, cmdKey) {
		cmd.Process.Signal(os.Interrupt)
		cmd.Wait()
		//cmd.Process.Kill()
		env = getEnv(client, envDir)
		command = getCmd(client, cmdKey)
		cmd = spawnProc(env, command)
	}
}
