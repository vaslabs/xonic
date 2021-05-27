package main

import (
	"os"

	"github.com/vaslabs/server"
	"github.com/vaslabs/client"

)

func main() {
	if (len(os.Args) > 1 && os.Args[1] == "server") {
		server.Run()
	} else { 
		client.Run()
	}
}