package main

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
