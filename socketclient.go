package main

import (
	"bufio"
	"fmt"
	"net"
)

type socketClient struct {
	chanInCmd   chan string
	chanOutCmd  chan string
	chanDispose chan bool
	reader      *bufio.Reader
	writer      *bufio.Writer
	conn        net.Conn
}

func (client *socketClient) read() {
	defer client.close()
	for {
		line, error := client.reader.ReadString('\r')
		if error != nil {
			fmt.Println(error)
			//client.chanDispose <- client
			break
		}
		client.chanInCmd <- line
	}
}

func (client *socketClient) write() {
	defer client.close()
	for {
		jsonString := <-client.chanOutCmd + "\r"
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

func (client *socketClient) open() {
	go client.read()
	go client.write()
}

func newSocketClient(_connection net.Conn) *socketClient {
	writer := bufio.NewWriter(_connection)
	reader := bufio.NewReader(_connection)

	client := &socketClient{
		chanInCmd:  make(chan string),
		chanOutCmd: make(chan string),
		//chanDispose: make(chan bool),
		reader: reader,
		writer: writer,
		conn:   _connection,
	}

	client.open()
	return client
}

func (client *socketClient) getInCmdChan() chan string {
	return client.chanInCmd
}

func (client *socketClient) getOutCmdChan() chan string {
	return client.chanOutCmd
}

func (client *socketClient) getType() int32 {
	return 0
}

func (client *socketClient) close() {
	client.conn.Close()
}
