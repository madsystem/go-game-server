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

type gameEntity struct {
	Pos       [2]float32 `json:"pos"`
	TargetPos [2]float32 `json:"targetPos"`
	Type      int32      `json:"entityType"`
	ID        int32      `json:"id"`
	Color     [3]uint32  `json:"color"`
	Velocity  float32    `json:"velocity"`
	Score     uint32     `json:"score"`

	lastUpdate time.Time

	chanInAction  chan string
	chanOutAction chan string
	chanAttack    chan attackInfo
	chanDead      chan *gameEntity
}

func newGameEntity(id int32, _chanInAction chan string, _chanOutAction chan string, _chanAttack chan attackInfo, _type int32) *gameEntity {
	var mapSizeX float32 = 100.0
	var mapSizeY float32 = 100.0
	startPosition := [2]float32{-mapSizeX/2 + rand.Float32()*mapSizeX,
		-mapSizeY/2 + rand.Float32()*mapSizeY}
	var velocity float32
	if _type == 0 {
		velocity = 8
	} else {
		velocity = 8 - rand.Float32()*4.0 // speed between 4 - 8 ms/s
	}

	newGameEntity := &gameEntity{
		Pos: startPosition,
		Color: [3]uint32{
			50 + rand.Uint32()%100,
			50 + rand.Uint32()%100,
			50 + rand.Uint32()%100,
		},

		Type:          _type,
		ID:            id,
		TargetPos:     startPosition,
		lastUpdate:    time.Now(),
		chanInAction:  _chanInAction,
		chanOutAction: _chanOutAction,
		chanAttack:    _chanAttack,
		Velocity:      velocity,
		Score:         0,
	}
	return newGameEntity
}

func (gameEntity *gameEntity) updateEntity() {
	var posVec vec2.T = gameEntity.Pos
	var targetPosVec vec2.T = gameEntity.TargetPos
	toTarget := vec2.Sub(&targetPosVec, &posVec)

	if toTarget.Length() > 0.2 {
		toTarget.Normalize()

		elapsed := time.Since(gameEntity.lastUpdate).Seconds()
		gameEntity.Pos[0] += toTarget[0] * gameEntity.Velocity * float32(elapsed)
		gameEntity.Pos[1] += toTarget[1] * gameEntity.Velocity * float32(elapsed)
	}

	gameEntity.lastUpdate = time.Now()
}

func (gameEntity *gameEntity) listen() {
	for {
		incAction := <-gameEntity.chanInAction
		fmt.Println("Received command", incAction)

		decoder := json.NewDecoder(strings.NewReader(incAction))

		var cmd clientBaseCmd
		err := decoder.Decode(&cmd)
		if err != nil {
			log.Fatal(err)
		}

		if cmd.Cmd == "move" {
			var gotoCmd clientGotoPosCmd
			err = json.Unmarshal(cmd.Payload, &gotoCmd)
			if err != nil {
				log.Println(err)
				continue
			}
			gameEntity.TargetPos = gotoCmd.TargetPos
		} else if cmd.Cmd == "attack" {
			var attackCmd clientAttackCmd
			err = json.Unmarshal(cmd.Payload, &attackCmd)
			if err != nil {
				log.Println(err)
				continue
			}

			gameEntity.chanAttack <- attackInfo{gameEntity.ID, int32(attackCmd.AttackTarget)}

		}
	}
}

func (gameEntity *gameEntity) dealDamage() {
	// kill entity directly
	gameEntity.chanDead <- gameEntity
}
