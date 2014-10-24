package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"github.com/csquared/forego/Godeps/_workspace/src/github.com/subosito/gotenv"
)

var envEntryRegexp = regexp.MustCompile("^([A-Za-z_0-9]+)=(.*)$")

type Env map[string]string

func GetEnv(c *etcd.Client, name string) (env Env) {
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

func ReadEnv(filename string) (Env, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return make(Env), nil
	}
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	env := make(Env)
	for key, val := range gotenv.Parse(fd) {
		env[key] = val
	}
	return env, nil
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
