package main

import ("net")

type NetworkHandler struct {
	clients  []*Client
	joins    chan net.Conn
	incoming chan string
	outgoing chan string
}

func (networkHandler *NetworkHandler) Broadcast(data string) {
	for _, client := range networkHandler.clients {
		client.chanOutAction <- data
	}
}

func (networkHandler *NetworkHandler) Join(connection net.Conn) {
	client := NewClient(connection)

	networkHandler.clients = append(networkHandler.clients, client)
	go func() {
		for {
			networkHandler.incoming <- <- client.chanInAction
		}
	}()
}

func (networkHandler *NetworkHandler) Listen() {
	go func() {
		for {
			select {
			case data := <-networkHandler.incoming:
				networkHandler.Broadcast(data)
			case conn := <-networkHandler.joins:
				networkHandler.Join(conn)
			}
		}
	}()
}

func NewNetworkHandler() *NetworkHandler {
	networkHandler := &NetworkHandler{
		clients:  make([]*Client, 0),
		joins:    make(chan net.Conn),
		incoming: make(chan string),
		outgoing: make(chan string),
	}

	networkHandler.Listen()

	return networkHandler
}
