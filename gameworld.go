package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Attack struct {
	attackerId int32
	targetId   int32
}

type GameWorld struct {
	gameEntities     []*GameEntity
	socketHandler    *SocketHandler
	websocketHandler *WebsocketHandler
	aiHandler        *AIHandler

	chanAttack chan Attack
	idCounter  int32
}

func NewGameWorld() *GameWorld {
	newGameWorld := &GameWorld{
		gameEntities: make([]*GameEntity, 0),
		chanAttack:   make(chan Attack),
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

func (gameWorld *GameWorld) IsAttackable(id int32) bool {
	for _, gameEntity := range gameWorld.gameEntities {
		if gameEntity.Id == id && gameEntity.Type != 0 {
			return true
		}
	}

	return false
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
	gameWorld.socketHandler = NewSocketHandler(gameWorld)
	gameWorld.socketHandler.Start()

	gameWorld.websocketHandler = NewWebsocketHandler(gameWorld)
	gameWorld.websocketHandler.Start()

	gameWorld.aiHandler = NewAIHandler(gameWorld)
	gameWorld.aiHandler.Start()
	// update clients
	go gameWorld.Update()

	fmt.Println("Server started ...")
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
			gameEntity.chanOutAction <- string(jsonCmd)
		}
	}
}

func (gameWorld *GameWorld) AddScore(id int32) {
	for _, gameEntity := range gameWorld.gameEntities {
		if gameEntity.Id == id {
			gameEntity.Score++
			break
		}
	}
}

func (gameWorld *GameWorld) UpdateAttacks() {
	select {
	case attack := <-gameWorld.chanAttack:
		if gameWorld.IsAttackable(attack.targetId) {
			// entity got attacked, kill it
			gameWorld.AddScore(attack.attackerId)
			gameWorld.RemoveGameEntity(attack.targetId)
		}
	default:
	}
}
