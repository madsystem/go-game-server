package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// aiClient handles
type aiClient struct {
	chanInCmd  chan string
	chanOutCmd chan string
	targetID   int32
	targetPos  [2]float32
}

func newAIClient() *aiClient {
	newClient := &aiClient{
		chanInCmd:  make(chan string),
		chanOutCmd: make(chan string),
		targetID:   -1,
		targetPos:  [2]float32{0, 0},
	}

	newClient.open()
	return newClient
}

func (client *aiClient) findTargetID(worldState worldStateCmd) int32 {
	var possibleTargets []int32
	for _, gameClient := range worldState.GameEntities {
		if gameClient.Type == 0 {
			// todo random target
			possibleTargets = append(possibleTargets, gameClient.ID)
		}
	}

	if len(possibleTargets) > 0 {
		return possibleTargets[rand.Intn(len(possibleTargets))]
	}

	return -1
}

func (client *aiClient) findTargetPos(id int32, worldState worldStateCmd) [2]float32 {
	for _, gameClient := range worldState.GameEntities {
		if gameClient.ID == id {
			return gameClient.Pos
		}
	}

	fmt.Println("AIHandler::FindTargetPos() id not found:", id)
	client.targetID = -1
	return client.targetPos
}

func (client *aiClient) read() {
	defer client.close()
	for {
		outCmd := <-client.chanOutCmd
		var worldState worldStateCmd
		json.Unmarshal([]byte(outCmd), &worldState)

		if client.targetID == -1 {
			client.targetID = client.findTargetID(worldState)
		}

		if client.targetID != -1 {
			targetPos := client.findTargetPos(client.targetID, worldState)
			client.targetPos = targetPos
		}
	}
}

func (client *aiClient) write() {
	defer client.close()
	for {
		time.Sleep(1 * time.Second)

		moveCmd := &clientGotoPosCmd{
			TargetPos: client.targetPos,
		}

		jsonCmd, _ := json.Marshal(moveCmd)
		baseMoveCmd := &clientBaseCmd{
			Cmd:     "move",
			Payload: jsonCmd,
		}

		moveJSON, _ := json.Marshal(baseMoveCmd)
		client.chanInCmd <- string(moveJSON)
	}
}

func (client *aiClient) open() {
	go client.read()
	go client.write()
}

func (client *aiClient) close() {

}

func (client *aiClient) getInCmdChan() chan string {
	return client.chanInCmd
}

func (client *aiClient) getOutCmdChan() chan string {
	return client.chanOutCmd
}

func (client *aiClient) getType() int32 {
	return 1
}
