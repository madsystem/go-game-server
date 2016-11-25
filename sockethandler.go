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

	client := newSocketClient(connection)
	id := handler.gameWorld.createGameEntity(client)
	handler.sendHandshake(client, id)
}

func (handler *socketHandler) sendHandshake(client *socketClient, id int32) {
	handshakeCmd := &handshake{
		ID: id,
	}
	jsonCmd, _ := json.Marshal(handshakeCmd)
	jsonOutString := string(jsonCmd)
	client.chanOutCmd <- jsonOutString
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
