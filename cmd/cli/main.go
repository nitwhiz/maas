package main

import (
	"github.com/docker/docker/client"
)

func main() {
	docker, err := client.NewClientWithOpts()

	if err != nil {
		panic(err)
	}

	InitCommands(docker)
}
