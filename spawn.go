package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

// Spawns process and waits for
// process to die or ENV to change
func SpawnLoop(n *Neuron) {
	n.Spawn()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			n.Kill()
			os.Exit(0)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(n.Ttl) * time.Second)
			n.HeartBeat()
		}
	}()

	for {
		envChan := make(chan interface{}, 1)
		processKilled := make(chan interface{}, 1)

		go func() {
			log.Println("action=wait-env")
			envChan <- n.Watch()
		}()
		go func() {
			log.Println("action=wait-process")
			processKilled <- n.Wait()
		}()

	Loop:
		for {
			select {
			case <-envChan:
				log.Println("action=env-changed")
				n.Kill()
				break Loop
			case <-processKilled:
				log.Println("action=process-killed")
				time.Sleep(1 * time.Second)
				break Loop
			}
		}

		if !n.Restart {
			break
		}
		n.Wait()
		n.Spawn()
	}
}

func spawnProc(env Env, command string) *exec.Cmd {
	log.Println("action=spawn-proc")

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
