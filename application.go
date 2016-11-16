package main

import (
	"fmt"
	"time"
)

func main() {
	gameWorld := NewGameWorld()
	fmt.Println("Setting up server ...")
	gameWorld.Start()

	for {
		time.Sleep(40 * time.Millisecond)
		//gameWorld.UpdateClients()
	}
}
