package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

type Neuron struct {
	AppName string
	Env     Env
	EnvName string
	CmdName string
	state   string
	Command string
	ProcId  string
	Restart bool
	Ttl     uint64
	Cmd     *exec.Cmd
	Etcd    *etcd.Client
}

func (n *Neuron) Spawn() exec.Cmd {
	n.Reload()
	n.State("starting")
	n.Cmd = spawnProc(n.Env, n.Command)
	n.State("up")
	return *n.Cmd
}

func (n *Neuron) State(state string) {
	n.state = state
	fmt.Printf("action=transition state=%s\n", state)
	n.HeartBeat()
}

//blocking call for update
func (n *Neuron) Watch() bool {
	watch(n.Etcd, n.EnvDir(), n.CmdKey())
	return true
}

func (n *Neuron) HeartBeat() {
	if n.Ttl > 0 {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}

		stateDir := n.StateDir()
		_, _ = n.Etcd.SetDir(stateDir, n.Ttl*3)
		_, _ = n.Etcd.Set(stateDir+"/state", n.state, 0)
		_, _ = n.Etcd.Set(stateDir+"/command", n.Command, 0)
		_, _ = n.Etcd.Set(stateDir+"/cmd", n.CmdKey(), 0)
		_, _ = n.Etcd.Set(stateDir+"/env", n.EnvDir(), 0)
		_, _ = n.Etcd.Set(stateDir+"/hostname", hostname, 0)
	}
}

func (n *Neuron) EnvDir() string {
	return "/services/" + n.AppName + "/envs/" + n.EnvName
}

func (n *Neuron) CmdKey() string {
	t := []string{"services", n.AppName, "processes", n.CmdName}
	return "/" + strings.Join(t, "/")
}

func (n *Neuron) StateDir() string {
	t := []string{"services", n.AppName, "running", n.CmdName, n.ProcId}
	return "/" + strings.Join(t, "/")
}

func (n *Neuron) Reload() {
	n.Env = GetEnv(n.Etcd, n.EnvDir())
	n.Command = GetCmd(n.Etcd, n.CmdKey())
}

func (n *Neuron) Wait() error {
	err := n.Cmd.Wait()
	n.State("down")
	return err
}

func (n *Neuron) Kill() {
	n.State("killing")
	n.Cmd.Process.Signal(os.Interrupt)
	n.State("down")
}
