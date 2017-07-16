package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// channel for save message sending other
	forward chan []byte
	// channel for client intend to join room. use for add client to `clients`
	join chan *client
	// channel for client intend to leave room. use for remove client from `clients`
	leave chan *client
	// save all joined clients
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// join this room
			r.clients[client] = true
		case client := <-r.leave:
			// leave from room
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// send message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send a message
				default:
					// fail to sending message
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// for using websocket. Upgrade HTTP connection by websocket.Upgrader
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
