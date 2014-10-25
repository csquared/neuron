package main

import (
	"os/exec"

	"github.com/coreos/go-etcd/etcd"
)

type Neuron struct {
	AppName string
	EnvDir  string
	Env     Env
	CmdKey  string
	Command string
	Restart bool
	Cmd     *exec.Cmd
	Etcd    *etcd.Client
}

func (n *Neuron) Spawn() exec.Cmd {
	n.Cmd = spawnProc(n.Env, n.Command)
	return *n.Cmd
}

//blocking call for update
func (n *Neuron) Watch() bool {
	return watch(n.Etcd, n.EnvDir, n.CmdKey)
}

func (n *Neuron) Reload() {
	n.Env = GetEnv(n.Etcd, n.EnvDir)
	n.Command = GetCmd(n.Etcd, n.CmdKey)
}
