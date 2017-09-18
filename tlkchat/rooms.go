package tlkchat

import (
	"fmt"
	"io"
	"sync"
)

type room struct {
	name string
	// All the sent messages will pass through this channel
	Msgch   chan string
	clients map[chan<- string]struct{}
	Quit    chan struct{}
	sync.RWMutex
}

// CreateRoom ...
func CreateRoom(name string) *room {
	r := &room{
		name:  name,
		Msgch: make(chan string),
		// clients channels can only receive data
		clients: make(map[chan<- string]struct{}),
		Quit:    make(chan struct{}),
	}
	r.Run()
	return r
}

// Run method executes a rotine that keeps waiting for incoming messages
func (r *room) Run() {
	logger.Println("Starting chat room ", r.name)
	// Using a go routine avoid the Run() method from blocking
	go func() {
		// Whenerver r.Msgch receives a value, the message will be broadcasted
		for msg := range r.Msgch {
			r.broadCastMsg(msg)
		}
	}()
}

// Clients will comunicate via TCP.
func (r *room) AddClient(conn io.ReadWriteCloser) {
	r.Lock()
	writeChannel, done := StartClient(r.Msgch, conn, r.Quit)
	// Add a new client
	r.clients[writeChannel] = struct{}{}
	r.Unlock()

	go func() {
		<-done
		r.RemoveClient(writeChannel)
	}()
}

func (r *room) RemoveClient(writeChannel chan<- string) {
	r.Lock()
	// Delete the client
	close(writeChannel)
	delete(r.clients, writeChannel)
	r.Unlock()

	// Using 'default' option avoid the 'select' statement from blocking
	select {
	case <-r.Quit:
		// If a 'r.Quit' signal has been sent and there are no more clients, the 'r.Msgch' is closed
		if r.ClCount() == 0 {
			close(r.Msgch)
		}
	default:
	}
}

func (r *room) ClCount() int {
	return len(r.clients)
}

func (r *room) broadCastMsg(msg string) {
	fmt.Println("Received message: ", msg)
	// Locks to protect the reading of the map of clients.
	r.RLock()
	defer r.RUnlock()
	for wc := range r.clients {
		/*
			To eficiently handle all the clients, the proccess of sending messages is made using go routines.
			Passing the channel as a parameter to the func, allow us to create a copy and ensure the data,
			because everytime the code inside the for block be executed, c maight change its value.
		*/
		go func(wc chan<- string) {
			wc <- msg
		}(wc)
	}
}
