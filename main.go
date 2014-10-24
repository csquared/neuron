package main

import (
	"flag"
	"os"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

type NeuronConfig struct {
	Cmd string
	Env Env
}

func main() {
	var envDir string
	flag.StringVar(&envDir, "env", "default", "name of env dir")
	var cmdKey string
	flag.StringVar(&cmdKey, "cmd", "", "name of cmd key")
	var etcdUrl string
	flag.StringVar(&etcdUrl, "etcd", "http://localhost:4001", "url of etcd")
	var restart bool
	flag.BoolVar(&restart, "r", false, "restart instead of crashing")
	flag.Parse()

	client := etcd.NewClient([]string{etcdUrl})

	if flag.Arg(0) == "bootstrap" {
		Bootstrap(client, flag.Arg(1))
		return
	}

	if flag.Arg(0) == "import" {
		Import(client)
		return
	}

	if envDir == "" || cmdKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	envDir, cmdKey = resolveKeys(envDir, cmdKey)

	//TODO: build up a config object and just
	//			pass it in to the event loop
	env := GetEnv(client, envDir)
	command := GetCmd(client, cmdKey)
	cmd := spawnProc(env, command)
	for watch(client, envDir, cmdKey) {
		cmd.Process.Signal(os.Interrupt)
		go func() {
			time.Sleep(10 * time.Second)
			cmd.Process.Kill()
		}()
		cmd.Wait()

		if !restart {
			os.Exit(0)
		}
		env = GetEnv(client, envDir)
		command = GetCmd(client, cmdKey)
		cmd = spawnProc(env, command)
	}
}
