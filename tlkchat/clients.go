package tlkchat

import (
	"bufio"
	"io"
)

type client struct {
	*bufio.Reader
	*bufio.Writer
	wc chan string
}

//StartClient is a channel generator and 'conn' is a connection received from a TCP client.
func StartClient(msgCh chan string, conn io.ReadWriteCloser, quit chan struct{}) (chan<- string, chan struct{}) {
	c := new(client)
	// Reads from the connection
	c.Reader = bufio.NewReader(conn)
	// Writes to the connection
	c.Writer = bufio.NewWriter(conn)
	c.wc = make(chan string)
	done := make(chan struct{})

	// Set up the reader
	go func() {
		scanner := bufio.NewScanner(c.Reader)
		/*
			Scan() returns false if the 'scanner' stops.
			Whenerver the 'scanner' receives a message through the 'conn', it will, enter the loop.
		*/
		for scanner.Scan() {
			msgCh <- scanner.Text()
		}
		// Once the 'scanner' stops (connection error, for example), the 'done' channel will be signaled
		done <- struct{}{}
	}()

	// Set up the writer
	c.writerMonitor()

	/*
		Every client awaits for a signal to be sent to the room Quit(r.Quit),
		and when it happens, the connection is closed.
		If a client is disconnected but the the room still running, just the done
		will be signaled and the select statement will exit.
	*/
	go func() {
		select {
		case <-quit:
			conn.Close()
		case <-done:
		}
	}()

	return c.wc, done
}

/*
	writerMonitor awaits for any sent message to the client channel and then writes it to
	the buffer writer (stream) of the connection.
*/
func (c *client) writerMonitor() {
	go func() {
		for s := range c.wc {
			// WriteString is a method from the embedded type bufio.Writer
			c.WriteString(s + "\n")
			c.Flush()
		}
	}()
}
