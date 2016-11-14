package main

import (
	"fmt"
	"math/rand"
)

type GameEntity struct {
	Pos  [2]float32 `json:"pos"`
	Vel  [2]float32 `json:"vel"`
	Type int32      `json:"entityType"`
	Id   int32      `json:"id"`

	chanInAction  chan string `json:"-"`
	chanOutAction chan string `json:"-"`
}

func NewGameEntity(id int32) *GameEntity {
	var mapSizeX float32 = 100.0
	var mapSizeY float32 = 100.0
	//maxVelX := 5
	//maxVelY := 5

	newGameEntity := &GameEntity{
		Pos: [2]float32{-mapSizeX/2 + rand.Float32()*mapSizeX,
			-mapSizeY/2 + rand.Float32()*mapSizeY},
		Vel:           [2]float32{0, 0},
		Type:          0,
		Id:            id,
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
