package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type gameWorld struct {
	gameEntities     []*gameEntity
	socketHandler    *socketHandler
	websocketHandler *websocketHandler
	aiHandler        *aiHandler

	chanClientCmd chan clientCmd
	idCounter     int32
}

func newGameWorld() *gameWorld {
	newgameWorld := &gameWorld{
		gameEntities:  make([]*gameEntity, 0),
		chanClientCmd: make(chan clientCmd),
	}

	return newgameWorld
}

func (gameWorld *gameWorld) fetchNewEntityID() int32 {
	newID := gameWorld.idCounter
	gameWorld.idCounter++
	return newID
}

func (gameWorld *gameWorld) createGameEntity(client client) int32 {
	gameEntity := newGameEntity(client, gameWorld.chanClientCmd)
	gameWorld.gameEntities = append(gameWorld.gameEntities, gameEntity)

	log.Println("AddGameEntity(): ", gameEntity, "Entity Count: ", len(gameWorld.gameEntities))
	return gameEntity.ID
}

func (gameWorld *gameWorld) isAttackable(id int32) bool {
	for _, gameEntity := range gameWorld.gameEntities {
		if gameEntity.ID == id && gameEntity.Type != 0 {
			return true
		}
	}

	return false
}

func (gameWorld *gameWorld) countNonHumanEntities() uint32 {
	var count uint32
	for _, gameEntity := range gameWorld.gameEntities {
		if gameEntity.Type != 0 {
			count++
		}
	}
	return count
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

		// update game world
		gameWorld.updateClientCommands()
		for _, gameEntity := range gameWorld.gameEntities {
			gameEntity.updateEntity()
		}

		// send update to clients
		worldStateCmd := newWorldStateCmd(gameWorld.gameEntities)
		jsonCmd, _ := json.Marshal(worldStateCmd)
		// update entities
		for _, gameEntity := range gameWorld.gameEntities {
			select {
			case gameEntity.clientOutCmd <- string(jsonCmd):
			default:
			}

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

func (gameWorld *gameWorld) updateClientCommands() {
	// currently client is only send
	select {
	case clientCmd := <-gameWorld.chanClientCmd:
		if clientCmd.getCmdType() == "AttackCmd" {
			attack := clientCmd.(*attackInfo)
			if gameWorld.isAttackable(attack.targetID) {
				// entity got attacked, kill it
				gameWorld.addScore(attack.attackerID)
				gameWorld.removeGameEntity(attack.targetID)
			}
		}
	default:
	}
}
