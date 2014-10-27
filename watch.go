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
	go n.Etcd.Watch(envDir, 0, true, envChan, nil)
	go n.Etcd.Watch(cmdKey, 0, false, cmdChan, nil)
	for _, dendriteDir := range n.Dendrites {
		log.Printf("action=wait-change dir=%s\n", dendriteDir)
		go n.Etcd.Watch(dendriteDir, 0, true, dendriteChan, nil)
	}

	log.Println("action=wait-change")

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
			log.Printf("action=dendrite-changed cmd=\"%s\"\n", r.Node.Key)
		}
	}
	return true
}
