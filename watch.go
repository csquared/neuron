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
		}
	}
	return true
}
