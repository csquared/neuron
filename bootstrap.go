package main

import (
	"log"

	"github.com/coreos/go-etcd/etcd"
)

func bootstrap(c *etcd.Client, language string) {
	appName := appName()
	_, err := c.SetDir("/services/"+appName+"/processes", 0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.Set("/services/"+appName+"/processes/web", "bundle exec puma -p $PORT", 0)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.SetDir("/services/"+appName+"/envs/dev", 0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.Set("/services/"+appName+"/envs/dev/PORT", "5000", 0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.Set("/services/"+appName+"/envs/dev/RACK_ENV", "development", 0)
	if err != nil {
		log.Fatal(err)
	}
}
