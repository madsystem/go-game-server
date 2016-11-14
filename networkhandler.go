package main

import (
	"fmt"
	"net"
)

type NetworkHandler struct {
	clients           []*Client
	joins             chan net.Conn
	gameWorld         *GameWorld
	connectionCounter int32
}

func (networkHandler *NetworkHandler) Join(connection net.Conn) {
	fmt.Println("Recv connection", connection)

	// create game entity and register it. not nice but works for now
	client := NewClient(connection)
	networkHandler.connectionCounter++ // inc id
	gameEntity := NewGameEntity(networkHandler.connectionCounter)
	networkHandler.gameWorld.AddGameEntity(gameEntity)
	networkHandler.clients = append(networkHandler.clients, client)

	// setup distribution channels
	go func() {
		for {
			select {
			case inAction := <-client.chanInAction:
				fmt.Println("Recv:", inAction)
				gameEntity.chanInAction <- inAction
			case outAction := <-gameEntity.chanOutAction:
				fmt.Println("Recv:", outAction)
				client.chanOutAction <- outAction
			}
		}
	}()
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
		clients:           make([]*Client, 0),
		joins:             make(chan net.Conn),
		gameWorld:         gameWorld,
		connectionCounter: 0,
	}

	return networkHandler
}

func (networkHandler *NetworkHandler) Start() {
	fmt.Println("Waiting for players ...")
	networkHandler.Listen()
}
