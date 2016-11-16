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
	Pos        [2]float32 `json:"pos"`
	TargetPos  [2]float32 `json:"targetPos"`
	Type       int32      `json:"entityType"`
	Id         int32      `json:"id"`
	maxSpeed   float32
	lastUpdate time.Time

	chanInAction  chan string
	chanOutAction chan string
}

func NewGameEntity(id int32, _chanInAction chan string, _chanOutAction chan string) *GameEntity {
	var mapSizeX float32 = 100.0
	var mapSizeY float32 = 100.0
	//maxVelX := 5
	//maxVelY := 5

	newGameEntity := &GameEntity{
		Pos: [2]float32{-mapSizeX/2 + rand.Float32()*mapSizeX,
			-mapSizeY/2 + rand.Float32()*mapSizeY},
		Type:          0,
		Id:            id,
		TargetPos:     [2]float32{0, 0},
		maxSpeed:      8,
		lastUpdate:    time.Now(),
		chanInAction:  _chanInAction,
		chanOutAction: _chanOutAction,
	}
	return newGameEntity
}

func (gameEntity *GameEntity) UpdateEntity() {
	for {
		time.Sleep(40 * time.Millisecond)
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
}

func (gameEntity *GameEntity) Listen() {
	for {
		time.Sleep(40 * time.Millisecond)
		select {
		case incAction := <-gameEntity.chanInAction:
			fmt.Println("Received command", incAction)

			decoder := json.NewDecoder(strings.NewReader(incAction))

			var cmd ClientBaseCmd
			err := decoder.Decode(&cmd)
			if err != nil {
				log.Fatal(err)
			}

			if cmd.Cmd == "move" {
				var gotoCmd ClientGotoCmd
				err = json.Unmarshal(cmd.Payload, &gotoCmd)
				if err != nil {
					log.Fatal(err)
				}

				//fmt.Println(gotoCmd)
				gameEntity.TargetPos = gotoCmd.TargetPos

			}
		default:

		}
	}
}
