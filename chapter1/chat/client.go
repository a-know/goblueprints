package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	// channel for sending message
	send chan []byte
	room *room
}
