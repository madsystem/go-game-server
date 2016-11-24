package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type websocketHandler struct {
	clients   []*websocketClient
	joins     chan websocket.Conn
	gameWorld *gameWorld
	upgrader  websocket.Upgrader
}

func (handler *websocketHandler) join(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket Recv connection")

	// create game entity and register it. not nice but works for now (create factory later)
	id := handler.gameWorld.fetchNewEntityID()
	var conn *websocket.Conn

	protocol := websocket.Subprotocols(r)
	if len(protocol) > 0 {
		newHeader := http.Header{"Sec-Websocket-Protocol": {protocol[0]}}
		conn, _ = handler.upgrader.Upgrade(w, r, newHeader)
	} else {
		conn, _ = handler.upgrader.Upgrade(w, r, nil)
	}

	client := newWebsocketClient(id, conn)
	gameEntity := newGameEntity(id, client.chanInCmd, client.chanOutCmd, handler.gameWorld.chanAttack, 0)

	// send handshake
	handler.sendHandshake(client)

	// add
	handler.gameWorld.addGameEntity(gameEntity)
	handler.clients = append(handler.clients, client)

	// setup distribution channels
	go func() {
		for {
			client := <-client.chanDisconnected
			fmt.Println("Player Disconnected")
			handler.gameWorld.removeGameEntity(client.id)
			handler.removeClient(client.id)
		}
	}()
}

func (handler *websocketHandler) sendHandshake(client *websocketClient) {
	handshakeCmd := &handshake{
		ID: client.id,
	}
	jsonCmd, _ := json.Marshal(handshakeCmd)
	jsonOutString := string(jsonCmd)

	client.chanOutCmd <- jsonOutString
}

func newWebsocketHandler(gameWorld *gameWorld) *websocketHandler {
	handler := &websocketHandler{
		clients:   make([]*websocketClient, 0),
		gameWorld: gameWorld,
	}

	return handler
}

func (handler *websocketHandler) start() {
	fmt.Println("Waiting for players ...")
	go handler.setupWebSocket()

}

func (handler *websocketHandler) removeClient(id int32) {
	for i, client := range handler.clients {
		if client.id == id {
			client.conn.Close()
			handler.clients = append(handler.clients[:i], handler.clients[i+1:]...)
			break
		}
	}
}

// websocket stuff

func alwaysTrue(r *http.Request) bool {
	return true
}

func (handler *websocketHandler) setupWebSocket() {

	// web socket
	handler.upgrader.CheckOrigin = alwaysTrue
	flag.Parse()
	log.SetFlags(0)
	var addr = flag.String("addr", ":4446", " server web socket port")

	http.HandleFunc("/", handler.join)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
