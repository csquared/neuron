package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-etcd/etcd"
)

func Import(c *etcd.Client, procfile, envfile string) {
	appName := appName()

	_, err := c.SetDir("/services/"+appName+"/processes", 0)
	_, err = c.SetDir("/services/"+appName+"/running", 0)

	fmt.Printf("action=import procfile=%s envfile=%s\n", procfile, envfile)

	//add processes from procfile
	if _, err := os.Stat(procfile); err == nil {
		procfile, err := ReadProcfile(procfile)
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range procfile.Entries {
			fmt.Printf("action=import-procfile process=%s\n", entry.Name)
			_, err = c.Set("/services/"+appName+"/processes/"+entry.Name, entry.Command, 0)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	//let's import .env to dev
	_, err = c.SetDir("/services/"+appName+"/envs/dev", 0)
	_, err = c.Set("/services/"+appName+"/envs/dev/PORT", "5000", 0)
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
			_, err = c.Set("/services/"+appName+"/envs/dev/"+key, value, 0)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
