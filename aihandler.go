package main

type AIHandler struct {
	//aiClients []*AIClient
	gameWorld *GameWorld
}

func NewAIHandler(_gameWorld *GameWorld) *AIHandler {
	aiHandler := &AIHandler{
		gameWorld: _gameWorld,
	}

	return aiHandler
}
