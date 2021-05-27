package main

import (
	"log"
	"os"

	"github.com/vaslabs/client"
	"github.com/vaslabs/server"
)

func main() {
	if (len(os.Args) <= 1) {
		log.Fatalln("Pass server or address to connect to")
	}
	if (len(os.Args) > 1 && os.Args[1] == "server") {
		server.Run()
	} else { 
		address := os.Args[1]
		client.Run(address)
	}
}