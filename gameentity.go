package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"strings"

	"github.com/ungerik/go3d/vec2"
)

type GameEntity struct {
	Pos       [2]float32 `json:"pos"`
	TargetPos [2]float32 `json:"targetPos"`
	Type      int32      `json:"entityType"`
	Id        int32      `json:"id"`
	Color     [3]uint32  `json:"color"`

	maxSpeed   float32
	lastUpdate time.Time

	chanInAction  chan string
	chanOutAction chan string
	chanAttack    chan int32
}

func NewGameEntity(id int32, _chanInAction chan string, _chanOutAction chan string, _chanAttack chan int32, _type int32) *GameEntity {
	var mapSizeX float32 = 100.0
	var mapSizeY float32 = 100.0
	//maxVelX := 5
	//maxVelY := 5
	startPosition := [2]float32{-mapSizeX/2 + rand.Float32()*mapSizeX,
		-mapSizeY/2 + rand.Float32()*mapSizeY}

	newGameEntity := &GameEntity{
		Pos: startPosition,
		Color: [3]uint32{
			50 + rand.Uint32()%100,
			50 + rand.Uint32()%100,
			50 + rand.Uint32()%100,
		},

		Type:          _type,
		Id:            id,
		TargetPos:     startPosition,
		maxSpeed:      8,
		lastUpdate:    time.Now(),
		chanInAction:  _chanInAction,
		chanOutAction: _chanOutAction,
		chanAttack:    _chanAttack,
	}
	return newGameEntity
}

func (gameEntity *GameEntity) UpdateEntity() {
	//time.Sleep(40 * time.Millisecond)
	var posVec vec2.T = gameEntity.Pos
	var targetPosVec vec2.T = gameEntity.TargetPos
	toTarget := vec2.Sub(&targetPosVec, &posVec)

	if toTarget.Length() > 0.2 {
		toTarget.Normalize()

		elapsed := time.Since(gameEntity.lastUpdate).Seconds()
		//log.Println("posVec", posVec, "targetPosVec", targetPosVec, "toTarget", toTarget, "elapsed", elapsed)
		gameEntity.Pos[0] += toTarget[0] * gameEntity.maxSpeed * float32(elapsed)
		gameEntity.Pos[1] += toTarget[1] * gameEntity.maxSpeed * float32(elapsed)
		//log.Println(gameEntity.Id, gameEntity.Pos)
	}

	gameEntity.lastUpdate = time.Now()
}

func (gameEntity *GameEntity) Listen() {
	for {
		incAction := <-gameEntity.chanInAction
		fmt.Println("Received command", incAction)

		decoder := json.NewDecoder(strings.NewReader(incAction))

		var cmd ClientBaseCmd
		err := decoder.Decode(&cmd)
		if err != nil {
			log.Fatal(err)
		}

		if cmd.Cmd == "move" {
			var gotoCmd ClientGotoPosCmd
			err = json.Unmarshal(cmd.Payload, &gotoCmd)
			if err != nil {
				log.Println(err)
				continue
			}
			//fmt.Println(gotoCmd)
			gameEntity.TargetPos = gotoCmd.TargetPos
		} else if cmd.Cmd == "attack" {
			var attackCmd ClientAttackCmd
			err = json.Unmarshal(cmd.Payload, &attackCmd)
			if err != nil {
				log.Println(err)
				continue
			}

			gameEntity.chanAttack <- attackCmd.AttackTarget

		}
	}
}
