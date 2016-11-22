package main

import "fmt"
import "github.com/gorilla/websocket"

type WebsocketClient struct {
	chanInCmd        chan string
	chanOutCmd       chan string
	chanDisconnected chan *WebsocketClient
	id               int32
	conn             *websocket.Conn
}

func NewWebsocketClient(_id int32, _conn *websocket.Conn) *WebsocketClient {
	client := &WebsocketClient{
		chanInCmd:        make(chan string),
		chanOutCmd:       make(chan string),
		chanDisconnected: make(chan *WebsocketClient),
		id:               _id,
		conn:             _conn,
	}

	client.Listen()
	return client
}

func (client *WebsocketClient) Read() {
	for {
		mt, msg, error := client.conn.ReadMessage()

		fmt.Println("msg type ", mt)

		if error != nil {
			fmt.Println(error)
			client.chanDisconnected <- client
			break
		}
		msgStr := string(msg[:len(msg)])
		client.chanInCmd <- msgStr
	}
}

func (client *WebsocketClient) Write() {
	for {
		jsonString := <-client.chanOutCmd
		jsonString = jsonString + "\r"
		error := client.conn.WriteMessage(2, []byte(jsonString)) // use binary message type
		if error != nil {
			fmt.Println(error)
			client.chanDisconnected <- client
			// todo close + remove client
			break
		}
	}
}

func (client *WebsocketClient) Listen() {
	go client.Read()
	go client.Write()
}
