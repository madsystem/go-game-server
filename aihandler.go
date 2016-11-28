package main

import "time"

type aiHandler struct {
	gameWorld *gameWorld
}

var spawnFrequency time.Duration = 10
var aiCount uint32 = 10

func newAIHandler(_gameWorld *gameWorld) *aiHandler {
	aiHandler := &aiHandler{
		gameWorld: _gameWorld,
	}

	return aiHandler
}

func (handler *aiHandler) start() {
	go handler.spawn()
}

func (handler *aiHandler) spawn() {
	for {
		if handler.gameWorld.countNonHumanEntities() < aiCount {
			createGameEntityMsg := newCreateGameEntityMsg(newAIClient())
			handler.gameWorld.createGameEntity <- *createGameEntityMsg
			<-createGameEntityMsg.responseID
		}
		time.Sleep(spawnFrequency * time.Second)
	}
}
