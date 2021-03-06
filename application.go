package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	gameWorld := newGameWorld()
	fmt.Println("Setting up server ...")
	gameWorld.Start()
	wg.Wait()
}
