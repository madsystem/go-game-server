package main



type GameWorld struct {
	gameEntities []*GameEntity 
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