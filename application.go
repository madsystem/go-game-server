package main

import (
	"fmt"
)

func main() {
	gameWorld := NewGameWorld()
	fmt.Println("Setting up server ...")
	gameWorld.Start()

	for {
		gameWorld.UpdateClients()
	}
}
