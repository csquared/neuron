package main

import (
	"log"

	"github.com/coreos/go-etcd/etcd"
)

func (n *Neuron) Watch() bool {
	envDir := n.EnvDir()
	cmdKey := n.CmdKey()

	envChan := make(chan *etcd.Response, 1)
	cmdChan := make(chan *etcd.Response, 1)
	dendriteChan := make(chan *etcd.Response, 1)
	log.Printf("action=wait-change dir=%s\n", envDir)
	go n.Etcd.Watch(envDir, 0, true, envChan, nil)
	log.Printf("action=wait-change key=%s\n", cmdKey)
	go n.Etcd.Watch(cmdKey, 0, false, cmdChan, nil)
	for _, dendriteDir := range n.Dendrites {
		log.Printf("action=wait-change dir=%s\n", dendriteDir)
		go n.Etcd.Watch(dendriteDir, 0, true, dendriteChan, nil)
	}

WatchLoop:
	for {
		select {
		case r := <-envChan:
			log.Printf("action=env-changed key=%s\n", r.Node.Key)
			break WatchLoop
		case r := <-cmdChan:
			log.Printf("action=cmd-changed cmd=\"%s\"\n", r.Node.Value)
			break WatchLoop
		case r := <-dendriteChan:
			log.Printf("action=dendrite-changed key=\"%s\"\n", r.Node.Key)
		}
	}
	return true
}
