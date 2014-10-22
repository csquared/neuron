package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

var envEntryRegexp = regexp.MustCompile("^([A-Za-z_0-9]+)=(.*)$")

type Env map[string]string

func getEnv(c *etcd.Client, name string) (env Env) {
	log.Printf("action=get-env name=%s\n", name)
	resp, err := c.Get(name, false, true)
	if err != nil {
		log.Fatal(err)
	}

	env = make(Env, len(resp.Node.Nodes))
	for _, n := range resp.Node.Nodes {
		cutpoint := strings.LastIndex(n.Key, "/") + 1
		key := n.Key[cutpoint:]
		env[key] = n.Value
		log.Printf("%s: %s\n", key, n.Value)
	}

	return
}

//thanks ddollar!
func (e *Env) asArray() (env []string) {
	for _, pair := range os.Environ() {
		env = append(env, pair)
	}
	for name, val := range *e {
		env = append(env, fmt.Sprintf("%s=%s", name, val))
	}
	return
}
