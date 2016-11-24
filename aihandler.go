package main

import "time"

type aiHandler struct {
	gameWorld *gameWorld
	clients   []*aiClient
}

var spawnFrequency time.Duration = 10

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
		// create game entity and register it. not nice but works for now (create factory later)
		id := handler.gameWorld.fetchNewEntityID()
		client := newAIClient()
		gameEntity := newGameEntity(id, client.chanInCmd, client.chanOutCmd, handler.gameWorld.chanAttack, 1)
		handler.gameWorld.addGameEntity(gameEntity)
		handler.clients = append(handler.clients, client)

		time.Sleep(spawnFrequency * time.Second)
	}
}
