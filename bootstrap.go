package main

import "log"

func Bootstrap(n *Neuron) {
	c := n.Etcd
	appName := n.AppName

	_, err := c.SetDir("/services/"+appName+"/processes", 0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.Set("/services/"+appName+"/processes/web", "bin/web", 0)
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
