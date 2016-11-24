package main

import "encoding/json"

type worldStateCmd struct {
	Cmd          string        `json:"cmd"`
	ID           int32         `json:"id"`
	GameEntities []*gameEntity `json:"payload"`
}

type clientBaseCmd struct {
	Cmd     string          `json:"cmd"`
	Payload json.RawMessage `json:"payload"`
}

type clientGotoPosCmd struct {
	TargetPos [2]float32 `json:"pos"`
}

type clientAttackCmd struct {
	AttackTarget int32 `json:"id"`
}

type handshake struct {
	ID int32 `json:"id"`
}

func newWorldStateCmd(gameEntities []*gameEntity) *worldStateCmd {
	newCmd := &worldStateCmd{
		Cmd:          "worldState",
		ID:           0,
		GameEntities: gameEntities,
	}

	return newCmd
}
