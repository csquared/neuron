package main

import "os"

func main() {
	var n = &Neuron{}
	Config(n)

	n.Reload()
	cmd := n.Spawn()
	for n.Watch() {
		cmd.Process.Signal(os.Interrupt)

		/*
					go func() {
						time.Sleep(10 * time.Second)
			//			cmd.Process.Kill()
					}()
		*/
		cmd.Wait()

		if !n.Restart {
			os.Exit(0)
		}
		n.Reload()
		cmd = n.Spawn()
	}
}
