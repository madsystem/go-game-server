package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type socketHandler struct {
	clients   []*socketClient
	joins     chan net.Conn
	gameWorld *gameWorld
}

func (handler *socketHandler) join(connection net.Conn) {
	fmt.Println("Recv connection", connection)

	// create game entity and register it. not nice but works for now (create factory later)
	id := handler.gameWorld.fetchNewEntityID()
	client := newSocketClient(connection, id)
	gameEntity := newGameEntity(id, client.chanInAction, client.chanOutAction, handler.gameWorld.chanAttack, 0)

	// send handshake
	handler.sendHandshake(client)

	// add to handler and to gameworld
	handler.gameWorld.addGameEntity(gameEntity)
	handler.clients = append(handler.clients, client)

	// setup distribution channels
	go func() {
		for {
			client := <-client.chanDisconnected
			fmt.Println("Player Disconnected")
			client.conn.Close()
			handler.gameWorld.removeGameEntity(client.id)
			handler.removeClient(client.id)
		}
	}()
}

func (handler *socketHandler) sendHandshake(client *socketClient) {
	handshakeCmd := &handshake{
		ID: client.id,
	}
	jsonCmd, _ := json.Marshal(handshakeCmd)
	jsonOutString := string(jsonCmd)

	client.chanOutAction <- jsonOutString
}

func (handler *socketHandler) listen() {
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
				handler.join(conn)
			}
		}
	}()

}

func newSocketHandler(gameWorld *gameWorld) *socketHandler {
	handler := &socketHandler{
		clients:   make([]*socketClient, 0),
		joins:     make(chan net.Conn),
		gameWorld: gameWorld,
	}

	return handler
}

func (handler *socketHandler) start() {
	fmt.Println("Waiting for players ...")
	handler.listen()

}

func (handler *socketHandler) removeClient(id int32) {
	for i, client := range handler.clients {
		if client.id == id {
			handler.clients = append(handler.clients[:i], handler.clients[i+1:]...)
			break
		}
	}
}
