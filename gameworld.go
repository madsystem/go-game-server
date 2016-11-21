package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type GameWorld struct {
	gameEntities     []*GameEntity
	socketHandler    *SocketHandler
	websocketHandler *WebsocketHandler

	chanAttack chan int32
	idCounter  int32
}

func NewGameWorld() *GameWorld {
	newGameWorld := &GameWorld{
		gameEntities: make([]*GameEntity, 0),
		chanAttack:   make(chan int32),
	}

	return newGameWorld
}

func (gameWorld *GameWorld) FetchNewEntityId() int32 {
	newId := gameWorld.idCounter
	gameWorld.idCounter++
	return newId
}

func (gameWorld *GameWorld) AddGameEntity(gameEntity *GameEntity) {
	gameWorld.gameEntities = append(gameWorld.gameEntities, gameEntity)

	// start entity loop
	go gameEntity.Listen()

	log.Println("AddGameEntity(): ", gameEntity, "Entity Count: ", len(gameWorld.gameEntities))
}

func (gameWorld *GameWorld) RemoveGameEntity(id int32) {
	log.Println("RemoveGameEntity: ", id)
	for i, gameEntity := range gameWorld.gameEntities {
		if gameEntity.Id == id {
			gameWorld.gameEntities = append(gameWorld.gameEntities[:i], gameWorld.gameEntities[i+1:]...)
			break
		}
	}
}

func (gameWorld *GameWorld) Start() {
	fmt.Println("Server started ...")
	socketHandler := NewSocketHandler(gameWorld)
	gameWorld.socketHandler = socketHandler
	socketHandler.Start()

	newWebSocketHandler := NewWebsocketHandler(gameWorld)
	gameWorld.websocketHandler = newWebSocketHandler
	newWebSocketHandler.Start()

	// update clients
	go gameWorld.Update()

	// create dummy monster
	gameWorld.AddGameEntity(NewGameEntity(gameWorld.FetchNewEntityId(), nil, nil, gameWorld.chanAttack, 1))
}

func (gameWorld *GameWorld) Update() {
	for {
		time.Sleep(40 * time.Millisecond) // sleep 40 ms

		gameWorld.UpdateAttacks()
		for _, gameEntity := range gameWorld.gameEntities {
			gameEntity.UpdateEntity()
		}

		// send update to clients
		updateWorldCmd := NewUpdateWorldStateCmd(gameWorld.gameEntities)
		jsonCmd, _ := json.Marshal(updateWorldCmd)
		// update entities
		for _, gameEntity := range gameWorld.gameEntities {
			if gameEntity.Type == 0 {
				gameEntity.chanOutAction <- string(jsonCmd)
				//fmt.Println("Update Client:", time.Now(), "Entity:", index, "OutString:", string(jsonOutString))
			}
		}
	}

}

func (gameWorld *GameWorld) UpdateAttacks() {
	// process attacks
	select {
	case attackTarget := <-gameWorld.chanAttack:
		gameWorld.RemoveGameEntity(attackTarget)
	default:
	}
}
