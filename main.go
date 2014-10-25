package main

import "fmt"

var banner string = `
   ____  ___  __  ___________  ____
  / __ \/ _ \/ / / / ___/ __ \/ __ \
 / / / /  __/ /_/ / /  / /_/ / / / /
/_/ /_/\___/\__,_/_/   \____/_/ /_/
`

func main() {
	fmt.Println(banner)
	var n = &Neuron{}
	Config(n)
	SpawnLoop(n)
}
