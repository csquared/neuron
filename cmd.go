package main

import (
	"log"

	"github.com/coreos/go-etcd/etcd"
)

func GetCmd(c *etcd.Client, name string) string {
	log.Printf("action=get-cmd name=%s\n", name)
	resp, err := c.Get(name, false, false)
	if err != nil {
		return name
	}

	return resp.Node.Value
}
