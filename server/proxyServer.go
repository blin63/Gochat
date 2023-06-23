/*
	proxyServer.go
	Name: Brendan Lin

	NOTE: Program execution for servers start here
*/

package main

import (
	"bufio"
	"fmt"
	"localhost/chat"
	"log"
	"net"
	"regexp"
	"strings"
)

// Caller functions to broadcast server; Validation is done in these functions
// Calls the broadcast server's list function
func list(conn net.Conn, c *chat.Chat) {
	c.List(conn)
}

// Calls the broadcast server's nick function
func nick(name string, conn net.Conn, c *chat.Chat) bool {
	// Check if nickname is valid
	match, _ := regexp.MatchString(`^[a-zA-Z][\w]{0,9}$`, name)

	// If nickname is invalid, return
	if !match {
		fmt.Fprintf(conn, "Invalid Nickname\n")
		return false
	}

	lst := c.ListInternal() // Get list of nicks

	// Check if nickname is already taken
	for _, v := range lst {
		if v == name {
			fmt.Fprintf(conn, "Nickname already taken\n")
			return false
		}
	}

	// Check if user already has a nickname;
	// If so, replace it with new nickname
	conns := c.ListConn()

	for i, v := range conns {
		if v == conn {
			if c.ListInternal()[i] != "" {
				c.UpdateNick(c.ListInternal()[i], name)
				conn.Write([]byte(fmt.Sprintf("Nickname changed to %s\n", name)))
				return true
			}
		}
	}

	// Call broadcast server's nick function
	c.Nick(name, conn)
	conn.Write([]byte(fmt.Sprintf("Nickname set to %s\n", name)))
	return true
}

// Calls the broadcast server's bc function
func bc(args []string, conn net.Conn, c *chat.Chat) bool {
	// Check if user has a nickname
	conns := c.ListConn()
	connFound := false

	for i, v := range conns {
		if v == conn {
			connFound = true
			if c.ListInternal()[i] == "" {
				fmt.Fprintf(conn, "You must have a nickname to broadcast\n")
				return false
			}
		}
	}

	// If connection is not found, return
	if !connFound {
		fmt.Fprintf(conn, "You must have a nickname to broadcast\n")
		return false
	}

	// Call broadcast server's Bc function
	createMsg := strings.Join(args, " ") // Create message from arguments
	c.Bc(createMsg)
	return true
}

// Handle Client Connection
func handle(conn net.Conn, c *chat.Chat) {
	defer conn.Close()

	// Go routine to check if client is still connected
	go func() {
		for {
			_, err := conn.Read([]byte{})
			if err != nil {
				// Print client disconnect on server
				fmt.Println("Client disconnected:", conn.RemoteAddr().String())

				// Get nickname of disconnected client
				disconnected_nick := c.GetNick(conn)

				// Remove client from chat server
				c.Remove(conn)

				// Send disconnect message to all clients
				c.Bc(fmt.Sprintf("%s has disconnected", disconnected_nick))
				return
			}
		}
	}()

	s := bufio.NewScanner(conn)
	for s.Scan() {
		// Split command and arguments to send to chat server
		ln := strings.Split(s.Text(), " ")
		cmd := ln[0]
		args := ln[1:]

		// Get command from client
		switch cmd {
		case "/LIST":
			list(conn, c)
		case "/NICK":
			// Check if nickname is provided
			if len(args) == 0 {
				fmt.Fprintf(conn, "No nickname provided\n")
				continue
			}

			nick_success := nick(args[0], conn, c)

			// Print nickname on server
			if nick_success {
				fmt.Println("Nickname OK")
			} else {
				fmt.Println("Nickname failed")
			}
		case "/BC":
			// Check if message is provided
			if len(args) == 0 {
				fmt.Fprintf(conn, "No message provided\n")
				continue
			}

			bc_success := bc(args, conn, c)

			// Print message on server
			if bc_success {
				fmt.Println("Message OK")
			} else {
				fmt.Println("Message failed")
			}
		default:
			fmt.Fprintf(conn, "Invalid Command\n")
		}
	}
}

// Start Proxy Server
func main() {
	// Listen for connections
	ln, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize chat server
	c := chat.ChatInit()

	// Handle connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		// Handle client connection
		go handle(conn, c)
	}
}
