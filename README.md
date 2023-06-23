Gochat

This is a simple demonstration chat server program written in Golang using TCP sockets.
The server must be built using go build proxyServer.go in the server folder first before the server can be run
Afterwards, the client can connect to the chat server using the client java program.

The below is the list of commands that can be used in the chat server:
/LIST - List all the users in the chat server
/NICK - Add/change the nickname of the user. Each user on server must have a unique nickname.
/BC - Broadcast a message to all the users in the chat server

NOTE: A user must have a nickname registered to the server before they can send any messages to other users
on the server. The nickname can be changed at any time using /NICK command.