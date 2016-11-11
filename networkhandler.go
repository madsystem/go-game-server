package main

import (
    "net"
    "fmt"
    )

type NetworkHandler struct {
	clients  []*Client
	joins    chan net.Conn
    gameWorld *GameWorld
}

func (networkHandler *NetworkHandler) Join(connection net.Conn) {
	client := NewClient(connection)
    gameEntity := NewGameEntity()

	networkHandler.clients = append(networkHandler.clients, client)

	go func() {
		for {
            // distribut data / actions
            gameEntity.chanInAction <- <- client.chanInAction
            //client.chanOutAction <- <- gameEntity.chanOutAction
            
		}
	}()
}

func (networkHandler *NetworkHandler) Listen() {
	go func() {
        listener,_ := net.Listen("tcp", ":4444")	
        for{
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
		clients:  make([]*Client, 0),
		joins:    make(chan net.Conn),
        gameWorld: gameWorld,
	}
	

	return networkHandler
}

func (networkHandler *NetworkHandler) Start(){
	fmt.Println("Waiting for players ...")
	networkHandler.Listen()
}
