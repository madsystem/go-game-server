package main

import "fmt"
import "github.com/gorilla/websocket"

type websocketClient struct {
	chanInCmd        chan string
	chanOutCmd       chan string
	chanDisconnected chan *websocketClient
	id               int32
	conn             *websocket.Conn
}

func newWebsocketClient(_id int32, _conn *websocket.Conn) *websocketClient {
	client := &websocketClient{
		chanInCmd:        make(chan string),
		chanOutCmd:       make(chan string),
		chanDisconnected: make(chan *websocketClient),
		id:               _id,
		conn:             _conn,
	}

	client.listen()
	return client
}

func (client *websocketClient) read() {
	defer client.conn.Close()
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

func (client *websocketClient) write() {
	defer client.conn.Close()
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

func (client *websocketClient) listen() {
	go client.read()
	go client.write()
}
