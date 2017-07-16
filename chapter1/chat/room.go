package main

type room struct {
	// channel for save message sending other
	forward chan []byte
}
