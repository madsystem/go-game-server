package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type GameWorld struct {
	gameEntities   []*GameEntity
	networkHandler *NetworkHandler

	testChan chan string
}

func NewGameWorld() *GameWorld {
	newGameWorld := &GameWorld{
		gameEntities: make([]*GameEntity, 0),
	}

	return newGameWorld
}

func (gameWorld *GameWorld) AddGameEntity(gameEntity *GameEntity) {
	gameWorld.gameEntities = append(gameWorld.gameEntities, gameEntity)
}

func (gameWorld *GameWorld) Start() {
	fmt.Println("Server started ...")
	newNetworkHandler := NewNetworkHandler(gameWorld)
	gameWorld.networkHandler = newNetworkHandler
	newNetworkHandler.Start()

	// update clients
	//go gameWorld.UpdateClients()
}

func (gameWorld *GameWorld) UpdateClients() {
	//fmt.Println(time.Now(), "Update World Start")
	time.Sleep(40) // sleep 40 ms
	updateWorldCmd := NewUpdateWorldStateCmd(gameWorld.gameEntities)
	jsonCmd, _ := json.Marshal(updateWorldCmd)
	jsonOutString := string(jsonCmd) + "\r"

	for index, gameEntity := range gameWorld.gameEntities {
		fmt.Println(time.Now(), "Entity:", index, "OutString:", string(jsonOutString))
		gameEntity.chanOutAction <- string(jsonOutString)
	}
	//fmt.Println(time.Now(), "Update World Done")

}
