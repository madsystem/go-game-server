package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	clients   []*WebsocketClient
	joins     chan websocket.Conn
	gameWorld *GameWorld
	upgrader  websocket.Upgrader
}

func (handler *WebsocketHandler) Join(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket Recv connection")

	// create game entity and register it. not nice but works for now (create factory later)
	id := handler.gameWorld.FetchNewEntityId()
	conn, _ := handler.upgrader.Upgrade(w, r, nil)
	client := NewWebsocketClient(id, conn)
	gameEntity := NewGameEntity(id, client.chanInCmd, client.chanOutCmd, handler.gameWorld.chanAttack, 0)

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

func (handler *WebsocketHandler) SendHandshake(client *WebsocketClient) {
	handshakeCmd := &Handshake{
		Id: client.id,
	}
	jsonCmd, _ := json.Marshal(handshakeCmd)
	jsonOutString := string(jsonCmd)

	client.chanOutCmd <- jsonOutString
}

func NewWebsocketHandler(gameWorld *GameWorld) *WebsocketHandler {
	handler := &WebsocketHandler{
		clients:   make([]*WebsocketClient, 0),
		gameWorld: gameWorld,
	}

	return handler
}

func (handler *WebsocketHandler) Start() {
	fmt.Println("Waiting for players ...")
	handler.SetupWebSocket()

}

func (handler *WebsocketHandler) RemoveClient(id int32) {
	for i, client := range handler.clients {
		if client.id == id {
			handler.clients = append(handler.clients[:i], handler.clients[i+1:]...)
			break
		}
	}
}

// websocket stuff

func alwaysTrue(r *http.Request) bool {
	return true
}

func (handler *WebsocketHandler) SetupWebSocket() {

	// web socket
	handler.upgrader.CheckOrigin = alwaysTrue
	flag.Parse()
	log.SetFlags(0)
	var addr = flag.String("addr", ":80", "web socket port")

	http.HandleFunc("/echo", handler.Join)
	http.HandleFunc("/ws", handler.Join)
	http.HandleFunc("/", handler.Join)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
