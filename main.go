package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/coreos/go-etcd/etcd"
)

func getServerList() []string {
	return []string{"http://127.0.0.1:4001/"}
}

func main() {
	client := etcd.NewClient(getServerList())

	var name string
	flag.StringVar(&name, "env", "default", "name under /env")
	var cmd string
	flag.StringVar(&cmd, "cmd", "", "name under /env")
	flag.Parse()

	envName := "/env/" + name
	env := getEnv(client, envName)

	/*
		if env["COMMAND_"+cmd] != "" {
			cmd = env["COMMAND_"+cmd]
		}
	*/

	spawnProc(env, cmd)
	for watchEnv(client, envName) {
		env = getEnv(client, envName)
		spawnProc(env, cmd)
	}
}

func spawnProc(env Env, command string) {
	log.Println("spawn Proc")

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Env = env.asArray()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stdout, stderr)

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
}
