package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-etcd/etcd"
)

func Import(c *etcd.Client) {
	appName := appName()

	_, err := c.SetDir("/services/"+appName+"/processes", 0)
	if err != nil {
		log.Fatal(err)
	}

	//add processes from procfile
	if _, err := os.Stat("Procfile"); err == nil {
		procfile, err := ReadProcfile("Procfile")
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range procfile.Entries {
			fmt.Println("Importing", entry.Name, "process from Procfile")
			_, err = c.Set("/services/"+appName+"/processes/"+entry.Name, entry.Command, 0)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	//let's import .env to dev
	_, err = c.SetDir("/services/"+appName+"/envs/dev", 0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.Set("/services/"+appName+"/envs/dev/PORT", "5000", 0)
	if err != nil {
		log.Fatal(err)
	}

	//add dev env from .env file
	if _, err := os.Stat(".env"); err == nil {
		env, err := ReadEnv(".env")
		if err != nil {
			log.Fatal(err)
		}

		for key, value := range env {
			fmt.Println("Importing", key, "from .env")
			_, err = c.Set("/services/"+appName+"/envs/dev/"+key, value, 0)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
