/*
	chat.go
	Name: Brendan Lin
*/

package chat

import (
	"fmt"
	"net"
)

// Chat struct contains neccessary info needed from each client
type Chat struct {
	// Slices in Chat are parallel slices
	nicks   []string   // slice of nicks
	sockets []net.Conn // slice of client connections
}

// Initialize Chat struct
func ChatInit() *Chat {
	c := new(Chat)
	c.nicks = make([]string, 0)
	c.sockets = make([]net.Conn, 0)
	return c
}

// List prints the elements in the clients slice
func (c *Chat) List(conn net.Conn) {
	// Print the nicks to the client
	conn.Write([]byte(fmt.Sprintf("%s", c.nicks) + "\n"))
}

// ListInternal returns the elements in the clients slice
func (c *Chat) ListInternal() []string {
	return c.nicks
}

// ListConn returns the elements in the clients slice
func (c *Chat) ListConn() []net.Conn {
	return c.sockets
}

// GetNick returns the nickname of a client
func (c *Chat) GetNick(conn net.Conn) string {
	for i, v := range c.sockets {
		if v == conn {
			return c.nicks[i]
		}
	}
	return ""
}

// UpdateNick updates the nickname of a client
func (c *Chat) UpdateNick(oldNick string, newNick string) {
	for i, v := range c.nicks {
		if v == oldNick {
			c.nicks[i] = newNick
		}
	}
}

// Nick adds a new client to the clients slice.
// This function assumes that the client meets requirements
func (c *Chat) Nick(nick string, conn net.Conn) {
	c.nicks = append(c.nicks, nick)
	c.sockets = append(c.sockets, conn)
}

// Bc broadcasts messages to all clients
func (c *Chat) Bc(msg string) {
	// Loop through all client sockets and send them the message
	send := msg + "\n"
	for _, conn := range c.sockets {
		conn.Write([]byte(send))
	}
}

// remove removes a connection and nickname from the chat slices
func (c *Chat) Remove(conn net.Conn) {
	for i, v := range c.sockets {
		if v == conn {
			c.sockets = append(c.sockets[:i], c.sockets[i+1:]...)
			c.nicks = append(c.nicks[:i], c.nicks[i+1:]...)
		}
	}
}
