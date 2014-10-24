package main

import (
	"log"
	"os"
	"os/exec"
)

func spawnProc(env Env, command string) *exec.Cmd {
	log.Println("spawn Proc")

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Env = env.asArray()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	return cmd
}
