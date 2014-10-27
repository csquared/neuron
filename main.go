package main

import (
	"flag"
	"fmt"
	"os"
)

var banner string = `
   ____  ___  __  ___________  ____
  / __ \/ _ \/ / / / ___/ __ \/ __ \
 / / / /  __/ /_/ / /  / /_/ / / / /
/_/ /_/\___/\__,_/_/   \____/_/ /_/
`

func main() {
	fmt.Println(banner)
	var n = &Neuron{}

	var procfile, envfile string
	//import flags
	flag.StringVar(&procfile, "p", "Procfile", "procfile location for import")
	flag.StringVar(&envfile, "e", ".env", ".env location for import")

	//calls flag.Parse
	Config(n)

	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "bootstrap":
			Bootstrap(n)
			return
		case "import":
			Import(n, procfile, envfile)
			return
		default:
			n.CmdName = args[0]
		}
	}

	if n.CmdName == "" {
		fmt.Println("You are missing the command: specify with -cmd\n")
		flag.Usage()
		os.Exit(1)
	}
	SpawnLoop(n)
}
