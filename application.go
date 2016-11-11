package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Server started ...")
	networkHandler := NewNetworkHandler()

	fmt.Println("Waiting for players ...")
	listener,_ := net.Listen("tcp", ":4444")	

	for{
		conn, _ := listener.Accept()
		networkHandler.joins <- conn
	}
}

