package main

import (
	"fmt"
	)

type GameWorld struct {
	gameEntities []*GameEntity 
	networkHandler *NetworkHandler
}

func NewGameWorld() *GameWorld {
	newGameWorld := &GameWorld{
		gameEntities:  make([]*GameEntity, 0),
	}

	return newGameWorld
}

func (gameWorld *GameWorld) AddGameEntity(gameEntity *GameEntity){
	gameWorld.gameEntities = append(gameWorld.gameEntities, gameEntity)
}

func (gameWorld *GameWorld) Start(){
	fmt.Println("Server started ...")
	newNetworkHandler := NewNetworkHandler(gameWorld)
	gameWorld.networkHandler = newNetworkHandler
	newNetworkHandler.Start()
	
}


