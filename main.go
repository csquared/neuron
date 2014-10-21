package main

import (
	"flag"

	"github.com/coreos/go-etcd/etcd"
)

func main() {
	var envDir string
	flag.StringVar(&envDir, "env", "default", "name of env dir")
	var cmdKey string
	flag.StringVar(&cmdKey, "cmd", "", "name of cmd key")
	var etcdUrl string
	flag.StringVar(&etcdUrl, "etcd", "http://127.0.0.1:4001", "url of etcd")
	flag.Parse()

	client := etcd.NewClient([]string{etcdUrl})
	envDir, cmdKey = resolveKeys(envDir, cmdKey)
	env := getEnv(client, envDir)
	command := getCmd(client, cmdKey)

	cmd := spawnProc(env, command)
	for watchEnv(client, envDir) {
		cmd.Process.Kill()
		env = getEnv(client, envDir)
		cmd = spawnProc(env, command)
	}
}
