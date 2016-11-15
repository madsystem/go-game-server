package main

import "encoding/json"

type UpdateWorldStateCmd struct {
	Cmd          string        `json:"cmd"`
	Id           int32         `json:"id"`
	GameEntities []*GameEntity `json:"payload"`
}

type ClientBaseCmd struct {
	Cmd     string          `json:"cmd"`
	Payload json.RawMessage `json:"payload"`
}

type ClientGotoCmd struct {
	TargetPos [2]float32 `json:"pos"`
}

func NewUpdateWorldStateCmd(gameEntities []*GameEntity) *UpdateWorldStateCmd {
	newCmd := &UpdateWorldStateCmd{
		Cmd:          "worldState",
		Id:           0,
		GameEntities: gameEntities,
	}

	return newCmd
}
