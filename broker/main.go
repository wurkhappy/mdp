package main

import (
	"github.com/wurkhappy/mdp"
	"os"
)

func main() {
	verbose := len(os.Args) >= 2 && os.Args[1] == "-v"
	broker := mdp.NewBroker("tcp://*:5555", verbose)
	defer broker.Close()

	broker.Run()
}
