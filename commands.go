package main

type UpdateWorldStateCmd struct {
	Cmd          string        `json:"cmd"`
	Id           int32         `json:"id"`
	GameEntities []*GameEntity `json:"payload"`
}

func NewUpdateWorldStateCmd(gameEntities []*GameEntity) *UpdateWorldStateCmd {
	newCmd := &UpdateWorldStateCmd{
		Cmd:          "worldState",
		Id:           0,
		GameEntities: gameEntities,
	}

	return newCmd
}
