package main

import (
	"fmt"
)

type GameEntity struct {
	pos        [2]float32
	vel        [2]float32

	chanInAction chan string 
	chanOutAction chan string 
}


func NewGameEntity() *GameEntity {
	newGameEntity := &GameEntity{
		pos : [2]float32{0,0},
		vel : [2]float32{0,0},
		chanInAction : make(chan string),
		chanOutAction : make(chan string),
	}
	go newGameEntity.Listen()
	return newGameEntity
}

func (gameEntity *GameEntity) Listen(){
	for {
		incAction := <-gameEntity.chanInAction
		fmt.Println("Received command", incAction)
			// do stuff
	}
}
