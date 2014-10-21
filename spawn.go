package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func spawnProc(env Env, command string) *exec.Cmd {
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
	return cmd
	/*
		log.Printf("Waiting for command to finish...")
		err = cmd.Wait()
		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	*/
}
