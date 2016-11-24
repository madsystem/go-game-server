package main

import (
	"bufio"
	"fmt"
	"net"
)

type socketClient struct {
	id               int32
	chanInAction     chan string
	chanOutAction    chan string
	chanDisconnected chan *socketClient
	reader           *bufio.Reader
	writer           *bufio.Writer
	conn             net.Conn
}

func (client *socketClient) read() {
	defer client.conn.Close()
	for {
		line, error := client.reader.ReadString('\r')
		if error != nil {
			fmt.Println(error)
			client.chanDisconnected <- client
			break
		}
		client.chanInAction <- line
	}
}

func (client *socketClient) write() {
	defer client.conn.Close()
	for {
		jsonString := <-client.chanOutAction + "\r"
		_, error := client.writer.WriteString(jsonString)
		if error != nil {
			fmt.Println(error)
			client.chanDisconnected <- client
			// todo close + remove client
			break
		}
		client.writer.Flush()
	}
}

func (client *socketClient) listen() {
	go client.read()
	go client.write()
}

func newSocketClient(_connection net.Conn, _id int32) *socketClient {
	writer := bufio.NewWriter(_connection)
	reader := bufio.NewReader(_connection)

	client := &socketClient{
		chanInAction:     make(chan string),
		chanOutAction:    make(chan string),
		chanDisconnected: make(chan *socketClient),
		id:               _id,
		reader:           reader,
		writer:           writer,
		conn:             _connection,
	}

	client.listen()
	return client
}
