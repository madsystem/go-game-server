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

	createGameEntity chan createGameEntityMsg
	updateWorld      chan bool
	chanClientCmd    chan clientCmd
	idCounter        int32
}

func newGameWorld() *gameWorld {
	newgameWorld := &gameWorld{
		gameEntities:     make([]*gameEntity, 0),
		chanClientCmd:    make(chan clientCmd),
		createGameEntity: make(chan createGameEntityMsg),
		updateWorld:      make(chan bool),
	}

	return newgameWorld
}

type createGameEntityMsg struct {
	responseID chan int32
	client     client
}

func newCreateGameEntityMsg(client client) *createGameEntityMsg {
	newCreateGameEntityMsg := &createGameEntityMsg{
		responseID: make(chan int32),
		client:     client,
	}

	return newCreateGameEntityMsg
}

func (gameWorld *gameWorld) handleWorld() {
	for {
		select {
		case createGameEntityMsg := <-gameWorld.createGameEntity:
			gameEntity := newGameEntity(createGameEntityMsg.client, gameWorld.chanClientCmd)
			gameWorld.gameEntities = append(gameWorld.gameEntities, gameEntity)
			createGameEntityMsg.responseID <- gameEntity.ID
			log.Println("AddGameEntity(): ", gameEntity, "Entity Count: ", len(gameWorld.gameEntities))
		case clientCmd := <-gameWorld.chanClientCmd:
			if clientCmd.getCmdType() == "AttackCmd" {
				attack := clientCmd.(*attackInfo)
				if gameWorld.isAttackable(attack.targetID) {
					// entity got attacked, kill it
					gameWorld.addScore(attack.attackerID)
					_, gameEntity := gameWorld.findEntityByID(attack.targetID)
					gameEntity.kill()
				}
			}
		case <-time.After(40 * time.Millisecond):
			//	cleanup
			var stillAlive []*gameEntity
			for _, gameEntity := range gameWorld.gameEntities {
				if gameEntity.isAlive() {
					stillAlive = append(stillAlive, gameEntity)
				}
			}
			gameWorld.gameEntities = stillAlive

			// send state to game entities and its connected client
			for _, gameEntity := range gameWorld.gameEntities {
				gameEntity.updateEntity()
			}

			worldStateCmd := newWorldStateCmd(gameWorld.gameEntities)
			jsonCmd, _ := json.Marshal(worldStateCmd)
			for _, gameEntity := range gameWorld.gameEntities {
				select {
				case gameEntity.clientOutCmd <- string(jsonCmd):
				default:
				}
			}
		}
	}
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

func (gameWorld *gameWorld) findEntityByID(id int32) (int32, *gameEntity) {
	for i, gameEntity := range gameWorld.gameEntities {
		if gameEntity.ID == id {
			return int32(i), gameEntity
		}
	}
	return -1, nil
}

func (gameWorld *gameWorld) Start() {
	gameWorld.socketHandler = newSocketHandler(gameWorld)
	gameWorld.socketHandler.start()

	gameWorld.websocketHandler = newWebsocketHandler(gameWorld)
	gameWorld.websocketHandler.start()

	gameWorld.aiHandler = newAIHandler(gameWorld)
	gameWorld.aiHandler.start()

	// update clients
	go gameWorld.handleWorld()

	fmt.Println("Server started ...")
}

func (gameWorld *gameWorld) addScore(id int32) {
	for _, gameEntity := range gameWorld.gameEntities {
		if gameEntity.ID == id {
			gameEntity.Score++
			break
		}
	}
}
