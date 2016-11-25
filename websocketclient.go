package main

import "fmt"
import "github.com/gorilla/websocket"

type websocketClient struct {
	chanInCmd        chan string
	chanOutCmd       chan string
	chanDisconnected chan *websocketClient
	conn             *websocket.Conn
	isAliveFlag      bool
}

func newWebsocketClient(_conn *websocket.Conn) *websocketClient {
	client := &websocketClient{
		chanInCmd:   make(chan string),
		chanOutCmd:  make(chan string),
		conn:        _conn,
		isAliveFlag: true,
	}

	client.listen()
	return client
}

func (client *websocketClient) read() {
	defer client.close()
	for {
		mt, msg, error := client.conn.ReadMessage()

		fmt.Println("msg type ", mt)

		if error != nil {
			fmt.Println(error)
			break
		}
		msgStr := string(msg[:len(msg)])
		client.chanInCmd <- msgStr
	}
}

func (client *websocketClient) write() {
	defer client.close()
	for {
		jsonString := <-client.chanOutCmd
		jsonString = jsonString + "\r"
		error := client.conn.WriteMessage(2, []byte(jsonString)) // use binary message type
		if error != nil {
			fmt.Println(error)
			break
		}
	}
}

func (client *websocketClient) listen() {
	go client.read()
	go client.write()
}

func (client *websocketClient) getInCmdChan() chan string {
	return client.chanInCmd
}

func (client *websocketClient) getOutCmdChan() chan string {
	return client.chanOutCmd
}

func (client *websocketClient) getType() int32 {
	return 0
}

func (client *websocketClient) isAlive() bool {
	return client.isAliveFlag
}

func (client *websocketClient) close() {
	client.isAliveFlag = false
	client.conn.Close()
}
