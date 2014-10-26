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

	size := len(resp.Node.Nodes)
	env = make(Env, size)
	keys := make([]string, size)
	for i, n := range resp.Node.Nodes {
		cutpoint := strings.LastIndex(n.Key, "/") + 1
		key := n.Key[cutpoint:]
		env[key] = n.Value
		keys[i] = key
	}

	log.Printf("action=get-env got-keys=%s\n", strings.Join(keys, ","))
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

func (e Env) Getenv(s string) string {
	return e[s]
}

/*
func (e Env) doSubstitutions(c *etcd.Client) {
	for key, val := range *e {
		if strings.HasPrefix(val, "neuron+") {
			tokens := strings.SplitAfterN(val, "://", 1)
			firstSlash := strings.Index(tokens[0], "/")
			name := tokens[0][:firstSlash]
			fmt.Println(name)
			e[key] = name

			//processUrl := strings.TrimPrefix(val, "neuron+")
			//processDir := "/services/" + appName() + "/running/"
			//resp, err := c.Get(name, false, false)
			//resp.Node.Value
		}
	}
}
*/
