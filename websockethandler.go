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
	joins     chan websocket.Conn
	gameWorld *gameWorld
	upgrader  websocket.Upgrader
}

func (handler *websocketHandler) join(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket Recv connection")

	var conn *websocket.Conn
	protocol := websocket.Subprotocols(r)
	if len(protocol) > 0 {
		newHeader := http.Header{"Sec-Websocket-Protocol": {protocol[0]}}
		conn, _ = handler.upgrader.Upgrade(w, r, newHeader)
	} else {
		conn, _ = handler.upgrader.Upgrade(w, r, nil)
	}

	client := newWebsocketClient(conn)
	id := handler.gameWorld.createGameEntity(client)
	handler.sendHandshake(client, id)
}

func (handler *websocketHandler) sendHandshake(client *websocketClient, id int32) {
	handshakeCmd := &handshake{
		ID: id,
	}
	jsonCmd, _ := json.Marshal(handshakeCmd)
	jsonOutString := string(jsonCmd)

	client.chanOutCmd <- jsonOutString
}

func newWebsocketHandler(gameWorld *gameWorld) *websocketHandler {
	handler := &websocketHandler{
		gameWorld: gameWorld,
	}

	return handler
}

func (handler *websocketHandler) start() {
	fmt.Println("Waiting for players ...")
	go handler.setupWebSocket()

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
