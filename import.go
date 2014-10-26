package main

import (
	"fmt"
	"log"
	"os"
)

func Import(n *Neuron, procfile, envfile string) {
	c := n.Etcd
	_, err := c.SetDir("/services/"+n.AppName+"/processes", 0)
	_, err = c.SetDir("/services/"+n.AppName+"/running", 0)

	fmt.Printf("action=import procfile=%s envfile=%s\n", procfile, envfile)

	//add processes from procfile
	if _, err := os.Stat(procfile); err == nil {
		procfile, err := ReadProcfile(procfile)
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range procfile.Entries {
			fmt.Printf("action=import-procfile process=%s\n", entry.Name)
			_, err = c.Set("/services/"+n.AppName+"/processes/"+entry.Name, entry.Command, 0)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	//let's import .env to dev
	_, err = c.SetDir("/services/"+n.AppName+"/envs/dev", 0)
	_, err = c.Set("/services/"+n.AppName+"/envs/dev/PORT", "5000", 0)
	if err != nil {
		log.Fatal(err)
	}

	//add dev env from .env file
	if _, err := os.Stat(envfile); err == nil {
		env, err := ReadEnv(envfile)
		if err != nil {
			log.Fatal(err)
		}

		for key, value := range env {
			fmt.Printf("action=import-env-var key=%s\n", key)
			_, err = c.Set("/services/"+n.AppName+"/envs/dev/"+key, value, 0)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
