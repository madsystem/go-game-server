package main

import "time"

type AIHandler struct {
	gameWorld *GameWorld
	clients   []*AIClient
}

var spawnFrequency time.Duration = 10

func NewAIHandler(_gameWorld *GameWorld) *AIHandler {
	aiHandler := &AIHandler{
		gameWorld: _gameWorld,
	}

	return aiHandler
}

func (handler *AIHandler) Start() {
	go handler.Spawn()
}

func (handler *AIHandler) Spawn() {
	for {
		// create game entity and register it. not nice but works for now (create factory later)
		id := handler.gameWorld.FetchNewEntityId()
		client := NewAIClient()
		gameEntity := NewGameEntity(id, client.chanInCmd, client.chanOutCmd, handler.gameWorld.chanAttack, 1)
		handler.gameWorld.AddGameEntity(gameEntity)
		handler.clients = append(handler.clients, client)

		time.Sleep(spawnFrequency * time.Second)
	}
}
