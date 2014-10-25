package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

// Spawns process and waits for
// process to die or ENV to change
func SpawnLoop(n *Neuron) {
	cmd := n.Spawn()

	for {
		envChan := make(chan interface{}, 1)
		processKilled := make(chan interface{}, 1)

		go func() {
			log.Println("wait for env change`")
			envChan <- n.Watch()
		}()
		go func() {
			log.Println("wait for process")
			processKilled <- cmd.Wait()
		}()

	Loop:
		for {
			select {
			case <-envChan:
				log.Println("Env changed")
				cmd.Process.Signal(os.Interrupt)
				break Loop
			case <-processKilled:
				log.Println("Process killed")
				time.Sleep(1 * time.Second)
				break Loop
			}
		}

		if !n.Restart {
			os.Exit(0)
		}
		cmd = n.Spawn()
	}
}

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
