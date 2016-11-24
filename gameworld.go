package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type attackInfo struct {
	attackerID int32
	targetID   int32
}

type gameWorld struct {
	gameEntities     []*gameEntity
	socketHandler    *socketHandler
	websocketHandler *websocketHandler
	aiHandler        *aiHandler

	chanAttack chan attackInfo
	idCounter  int32
}

func newGameWorld() *gameWorld {
	newgameWorld := &gameWorld{
		gameEntities: make([]*gameEntity, 0),
		chanAttack:   make(chan attackInfo),
	}

	return newgameWorld
}

func (gameWorld *gameWorld) fetchNewEntityID() int32 {
	newID := gameWorld.idCounter
	gameWorld.idCounter++
	return newID
}

func (gameWorld *gameWorld) addGameEntity(gameEntity *gameEntity) {
	gameWorld.gameEntities = append(gameWorld.gameEntities, gameEntity)

	// start entity loop
	go gameEntity.listen()
	log.Println("AddGameEntity(): ", gameEntity, "Entity Count: ", len(gameWorld.gameEntities))
}

func (gameWorld *gameWorld) isAttackable(id int32) bool {
	for _, gameEntity := range gameWorld.gameEntities {
		if gameEntity.ID == id && gameEntity.Type != 0 {
			return true
		}
	}

	return false
}

func (gameWorld *gameWorld) removeGameEntity(id int32) {
	log.Println("RemoveGameEntity: ", id)
	for i, gameEntity := range gameWorld.gameEntities {
		if gameEntity.ID == id {
			gameWorld.gameEntities = append(gameWorld.gameEntities[:i], gameWorld.gameEntities[i+1:]...)
			break
		}
	}
}

func (gameWorld *gameWorld) Start() {
	gameWorld.socketHandler = newSocketHandler(gameWorld)
	gameWorld.socketHandler.start()

	gameWorld.websocketHandler = newWebsocketHandler(gameWorld)
	gameWorld.websocketHandler.start()

	gameWorld.aiHandler = newAIHandler(gameWorld)
	gameWorld.aiHandler.start()
	// update clients
	go gameWorld.update()

	fmt.Println("Server started ...")
}

func (gameWorld *gameWorld) update() {
	for {
		time.Sleep(40 * time.Millisecond) // sleep 40 ms

		gameWorld.updateAttacks()
		for _, gameEntity := range gameWorld.gameEntities {
			gameEntity.updateEntity()
		}

		// send update to clients
		worldStateCmd := newWorldStateCmd(gameWorld.gameEntities)
		jsonCmd, _ := json.Marshal(worldStateCmd)
		// update entities
		for _, gameEntity := range gameWorld.gameEntities {
			gameEntity.chanOutAction <- string(jsonCmd)
		}
	}
}

func (gameWorld *gameWorld) addScore(id int32) {
	for _, gameEntity := range gameWorld.gameEntities {
		if gameEntity.ID == id {
			gameEntity.Score++
			break
		}
	}
}

func (gameWorld *gameWorld) updateAttacks() {
	select {
	case attack := <-gameWorld.chanAttack:
		if gameWorld.isAttackable(attack.targetID) {
			// entity got attacked, kill it
			gameWorld.addScore(attack.attackerID)
			gameWorld.removeGameEntity(attack.targetID)
		}
	default:
	}
}
