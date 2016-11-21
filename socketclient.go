package main

import (
	"bufio"
	"fmt"
	"net"
)

type SocketClient struct {
	id               int32
	chanInAction     chan string
	chanOutAction    chan string
	chanDisconnected chan *SocketClient
	reader           *bufio.Reader
	writer           *bufio.Writer
}

func (client *SocketClient) Read() {
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

func (client *SocketClient) Write() {
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

func (client *SocketClient) Listen() {
	go client.Read()
	go client.Write()
}

func NewSocketClient(connection net.Conn, _id int32) *SocketClient {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &SocketClient{
		chanInAction:     make(chan string),
		chanOutAction:    make(chan string),
		chanDisconnected: make(chan *SocketClient),
		id:               _id,
		reader:           reader,
		writer:           writer,
	}

	client.Listen()

	return client
}
