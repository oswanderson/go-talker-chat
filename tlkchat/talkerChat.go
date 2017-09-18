// Talker Chat ...

package tlkchat

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/wand-go-talker-chat/ctmlogger"
)

var logger = ctmlogger.GetInstance()

// Run function starts a TCP server
func Run(connection string) error {
	listener, err := net.Listen("tcp", connection)
	if err != nil {
		logger.Println("Error at trying to create s char server:", err)
		return err
	}
	r := CreateRoom("TalkerChatRoom")

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch

		listener.Close()
		fmt.Println("Closing the char server...")
		close(r.Quit)

		/*
			Block the code until every client get closed.
			See the select statement in RemoveClient inside rooms.go
		*/
		if r.ClCount() == 0 {
			<-r.Msgch
		}
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Println("Error at accepting incoming connection:", err)
			break
		}
		handleConnection(r, conn)
	}

	return nil
}

func handleConnection(r *room, conn net.Conn) {
	logger.Println("Received request from client:", conn.RemoteAddr())
	r.AddClient(conn)
}
