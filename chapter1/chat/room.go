package main

import (
	"log"
	"net/http"
	"os"

	"github.com/a-know/goblueprints/chapter1/trace"
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
	// logging
	tracer trace.Tracer
}

func newRoom(logging bool) *room {
	tracer := trace.New(os.Stdout)
	if !logging {
		tracer = trace.Off()
	}
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  tracer,
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// join this room
			r.clients[client] = true
			r.tracer.Trace("Joined a new client.")
		case client := <-r.leave:
			// leave from room
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Leave a client.")
		case msg := <-r.forward:
			r.tracer.Trace("Receive a message: ", string(msg))
			// send message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send a message
					r.tracer.Trace(" -- Send a message: ", string(msg))
				default:
					// fail to sending message
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- Failed to send a message. Cleanup client...")
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
