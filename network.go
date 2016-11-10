package main

import (
	"fmt"
	"net"
	"encoding/json"
)

type JsonCommand struct{
	id string
	json string
}

type NetworkHandle struct {

}

func handleConnections() {
	service := ":4444" // the devils port :D
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// run as a goroutine
		go handleClientCommands(conn)
	}
}


func handleClientCommands(conn net.Conn) {
	fmt.Println("handle client commands")

	// close connection on exit
	defer conn.Close()
	//b := bufio.NewReader(conn)
	decoder := json.NewDecoder(conn) // maybe this works for json
	for {
		var msg JsonCommand
		err := decoder.Decode(&msg)
		if err != nil { // EOF, or worse
            break
        }
		fmt.Println(msg, err)
	}
}

func handleClientJoinedCmd() {
}

func clientActionCmd() {
}
