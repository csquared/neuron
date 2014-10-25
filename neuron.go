package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/coreos/go-etcd/etcd"
)

type Neuron struct {
	AppName  string
	EnvDir   string
	Env      Env
	CmdKey   string
	StateDir string
	state    string
	Command  string
	ProcId   string
	Restart  bool
	Ttl      uint64
	Cmd      *exec.Cmd
	Etcd     *etcd.Client
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
	watch(n.Etcd, n.EnvDir, n.CmdKey)
	return true
}

func (n *Neuron) HeartBeat() {
	if n.Ttl > 0 {
		n.StateDir = "/services/" + n.AppName + "/running/" + n.ProcId
		_, _ = n.Etcd.SetDir("/services/"+n.AppName+"/running", 0)
		_, _ = n.Etcd.SetDir(n.StateDir, n.Ttl*2)
		_, _ = n.Etcd.Set(n.StateDir+"/command", n.Command, 0)
		_, _ = n.Etcd.Set(n.StateDir+"/cmd", n.CmdKey, 0)
		_, _ = n.Etcd.Set(n.StateDir+"/env", n.EnvDir, 0)
		_, _ = n.Etcd.Set(n.StateDir+"/state", n.state, 0)
	}
}

func (n *Neuron) Reload() {
	n.Env = GetEnv(n.Etcd, n.EnvDir)
	n.Command = GetCmd(n.Etcd, n.CmdKey)
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
