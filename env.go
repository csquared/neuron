package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
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
	keys := make([]string, size)
	env = make(Env, size)
	for i, n := range resp.Node.Nodes {
		key := path.Base(n.Key)
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

func (e *Env) doSubstitutions(c *etcd.Client) (env Env, subs []string) {
	env = make(Env, len(*e))
	for key, val := range *e {
		env[key] = val
		if strings.HasPrefix(val, "neuron+") {
			raw_url := strings.TrimPrefix(val, "neuron+")
			url, err := url.Parse(raw_url)
			if err != nil {
				log.Println(err)
				continue
			}

			names := strings.Split(url.Host, ":")
			var tokens []string
			if len(names) == 1 {
				tokens = []string{"services", GetAppName(), "running", names[0]}
			} else {
				tokens = []string{"services", names[0], "running", names[1]}
			}
			dir := "/" + strings.Join(tokens, "/")

			subs = append(subs, dir)
			resp, err := c.Get(dir, true, true)
			if err != nil {
				log.Println(err)
				continue
			}

			n := len(resp.Node.Nodes)
			if n != 0 {
				i := 0
				live := resp.Node.Nodes[i]

				var hostname, port string
				for _, node := range live.Nodes {
					keyName := path.Base(node.Key)
					switch keyName {
					case "hostname":
						hostname = node.Value
					case "port":
						port = node.Value
					}
				}

				url.Host = hostname + ":" + port
				fmt.Println(val)
				env[key] = url.String()
			}
		}
	}

	return
}
