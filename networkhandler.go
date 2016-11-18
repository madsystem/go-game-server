package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type NetworkHandler struct {
	clients   []*Client
	joins     chan net.Conn
	gameWorld *GameWorld
}

func (networkHandler *NetworkHandler) Join(connection net.Conn) {
	fmt.Println("Recv connection", connection)

	// create game entity and register it. not nice but works for now (create factory later)
	id := networkHandler.gameWorld.FetchNewEntityId()
	client := NewClient(connection, id)
	gameEntity := NewGameEntity(id, client.chanInAction, client.chanOutAction, networkHandler.gameWorld.chanAttack, 0)

	// send handshake
	networkHandler.SendHandshake(client)

	// add
	networkHandler.gameWorld.AddGameEntity(gameEntity)
	networkHandler.clients = append(networkHandler.clients, client)

	// setup distribution channels
	go func() {
		for {
			client := <-client.chanDisconnected
			fmt.Println("Player Disconnected")
			networkHandler.gameWorld.RemoveGameEntity(client.id)
			networkHandler.RemoveClient(client.id)
		}
	}()
}

func (networkHandler *NetworkHandler) SendHandshake(client *Client) {
	handshakeCmd := &Handshake{
		Id: client.id,
	}
	jsonCmd, _ := json.Marshal(handshakeCmd)
	jsonOutString := string(jsonCmd) + "\r"

	client.chanOutAction <- jsonOutString
}

func (networkHandler *NetworkHandler) Listen() {
	go func() {
		listener, _ := net.Listen("tcp", ":4444")
		for {
			conn, _ := listener.Accept()
			networkHandler.joins <- conn
		}
	}()
	go func() {
		for {
			select {
			case conn := <-networkHandler.joins:
				networkHandler.Join(conn)
			}
		}
	}()
}

func NewNetworkHandler(gameWorld *GameWorld) *NetworkHandler {
	networkHandler := &NetworkHandler{
		clients:   make([]*Client, 0),
		joins:     make(chan net.Conn),
		gameWorld: gameWorld,
	}

	return networkHandler
}

func (networkHandler *NetworkHandler) Start() {
	fmt.Println("Waiting for players ...")
	networkHandler.Listen()
}

func (networkHandler *NetworkHandler) RemoveClient(id int32) {
	for i, client := range networkHandler.clients {
		if client.id == id {
			networkHandler.clients = append(networkHandler.clients[:i], networkHandler.clients[i+1:]...)
			break
		}
	}
}
