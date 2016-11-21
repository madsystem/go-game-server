package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type SocketHandler struct {
	clients   []*SocketClient
	joins     chan net.Conn
	gameWorld *GameWorld
}

func (handler *SocketHandler) Join(connection net.Conn) {
	fmt.Println("Recv connection", connection)

	// create game entity and register it. not nice but works for now (create factory later)
	id := handler.gameWorld.FetchNewEntityId()
	client := NewSocketClient(connection, id)
	gameEntity := NewGameEntity(id, client.chanInAction, client.chanOutAction, handler.gameWorld.chanAttack, 0)

	// send handshake
	handler.SendHandshake(client)

	// add
	handler.gameWorld.AddGameEntity(gameEntity)
	handler.clients = append(handler.clients, client)

	// setup distribution channels
	go func() {
		for {
			client := <-client.chanDisconnected
			fmt.Println("Player Disconnected")
			handler.gameWorld.RemoveGameEntity(client.id)
			handler.RemoveClient(client.id)
		}
	}()
}

func (handler *SocketHandler) SendHandshake(client *SocketClient) {
	handshakeCmd := &Handshake{
		Id: client.id,
	}
	jsonCmd, _ := json.Marshal(handshakeCmd)
	jsonOutString := string(jsonCmd)

	client.chanOutAction <- jsonOutString
}

func (handler *SocketHandler) Listen() {
	go func() {
		listener, _ := net.Listen("tcp", ":4444")
		for {
			conn, _ := listener.Accept()
			handler.joins <- conn
		}
	}()

	go func() {
		for {
			select {
			case conn := <-handler.joins:
				handler.Join(conn)
			}
		}
	}()

}

func NewSocketHandler(gameWorld *GameWorld) *SocketHandler {
	handler := &SocketHandler{
		clients:   make([]*SocketClient, 0),
		joins:     make(chan net.Conn),
		gameWorld: gameWorld,
	}

	return handler
}

func (handler *SocketHandler) Start() {
	fmt.Println("Waiting for players ...")
	handler.Listen()

}

func (handler *SocketHandler) RemoveClient(id int32) {
	for i, client := range handler.clients {
		if client.id == id {
			handler.clients = append(handler.clients[:i], handler.clients[i+1:]...)
			break
		}
	}
}
