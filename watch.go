package main

import (
	"log"

	"github.com/coreos/go-etcd/etcd"
)

func watch(c *etcd.Client, envDir, cmdKey string) bool {
	envChan := make(chan *etcd.Response)
	go c.Watch(envDir, 0, true, envChan, nil)

	cmdChan := make(chan *etcd.Response)
	go c.Watch(cmdKey, 0, false, cmdChan, nil)

	log.Println("Waiting for an update...")

	selected := false
	for !selected {
		select {
		case r := <-envChan:
			log.Printf("Got updated env: %s\n", r.Node.Key)
			selected = true
		case r := <-cmdChan:
			log.Printf("Got new command: %s\n", r.Node.Value)
			selected = true
		}
	}
	return true
}
