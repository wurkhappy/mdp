package main

import (
	//"flag"
	//"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/mdp"
	"os"
)

//var production = flag.Bool("production", false, "Production settings")

func main() {
	//flag.Parse()
	//if *production {
	//	config.Prod()
	//} else {
		//config.Test()
	//}
	verbose := len(os.Args) >= 2 && os.Args[1] == "-v"
	broker := mdp.NewBroker("tcp://*:5555", verbose)
	defer broker.Close()

	broker.Run()
}
