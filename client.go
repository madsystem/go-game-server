package main

type client interface {
	getType() int32
	getInCmdChan() chan string
	getOutCmdChan() chan string
}
