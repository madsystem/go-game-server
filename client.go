package main

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	id               int32
	chanInAction     chan string
	chanOutAction    chan string
	chanDisconnected chan *Client
	reader           *bufio.Reader
	writer           *bufio.Writer
}

func (client *Client) Read() {
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

func (client *Client) Write() {
	for {
		jsonString := <-client.chanOutAction
		_, error := client.writer.WriteString(jsonString)
		if error != nil {
			fmt.Println(error)
			//client.chanDisconnected <- client
			// todo close + remove client
			break
		}
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn, _id int32) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		chanInAction:     make(chan string),
		chanOutAction:    make(chan string),
		chanDisconnected: make(chan *Client),
		id:               _id,
		reader:           reader,
		writer:           writer,
	}

	client.Listen()

	return client
}
