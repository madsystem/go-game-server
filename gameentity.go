package main

import (
	"fmt"
)

type GameEntity struct {
	Pos  [2]float32 `json:"pos"`
	Vel  [2]float32 `json:"vel"`
	Type int        `json:"entityType"`
	Id   int        `json:"id"`

	chanInAction  chan string `json:"-"`
	chanOutAction chan string `json:"-"`
}

func NewGameEntity() *GameEntity {
	newGameEntity := &GameEntity{
		Pos:           [2]float32{0, 0},
		Vel:           [2]float32{0, 0},
		Type:          0,
		Id:            0,
		chanInAction:  make(chan string),
		chanOutAction: make(chan string),
	}
	go newGameEntity.Listen()
	return newGameEntity
}

func (gameEntity *GameEntity) Listen() {
	for {
		incAction := <-gameEntity.chanInAction
		fmt.Println("Received command", incAction)
		// do stuff
	}
}
