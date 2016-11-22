package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type AIClient struct {
	chanInCmd  chan string
	chanOutCmd chan string
	targetId   int32
	targetPos  [2]float32
}

func NewAIClient() *AIClient {
	newClient := &AIClient{
		chanInCmd:  make(chan string),
		chanOutCmd: make(chan string),
		targetId:   -1,
		targetPos:  [2]float32{0, 0},
	}

	newClient.Listen()
	return newClient
}

func (client *AIClient) FindTargetId(worldState UpdateWorldStateCmd) int32 {
	var possibleTargets []int32
	for _, gameClient := range worldState.GameEntities {
		if gameClient.Type == 0 {
			// todo random target
			possibleTargets = append(possibleTargets, gameClient.Id)
		}
	}

	if len(possibleTargets) > 0 {
		return possibleTargets[rand.Intn(len(possibleTargets))]
	}

	return -1
}

func (client *AIClient) FindTargetPos(id int32, worldState UpdateWorldStateCmd) [2]float32 {
	for _, gameClient := range worldState.GameEntities {
		if gameClient.Id == id {
			return gameClient.Pos
		}
	}

	fmt.Println("AIHandler::FindTargetPos() id not found:", id)
	client.targetId = -1
	return client.targetPos
}

func (client *AIClient) Read() {
	for {
		outCmd := <-client.chanOutCmd
		var worldState UpdateWorldStateCmd
		json.Unmarshal([]byte(outCmd), &worldState)

		if client.targetId == -1 {
			client.targetId = client.FindTargetId(worldState)
		}

		if client.targetId != -1 {
			targetPos := client.FindTargetPos(client.targetId, worldState)
			client.targetPos = targetPos
		}
	}
}

func (client *AIClient) Write() {
	for {
		time.Sleep(1 * time.Second)

		moveCmd := &ClientGotoPosCmd{
			TargetPos: client.targetPos,
		}

		jsonCmd, _ := json.Marshal(moveCmd)
		baseMoveCmd := &ClientBaseCmd{
			Cmd:     "move",
			Payload: jsonCmd,
		}

		moveJson, _ := json.Marshal(baseMoveCmd)
		client.chanInCmd <- string(moveJson)
	}
}

func (client *AIClient) Listen() {
	go client.Read()
	go client.Write()
}
