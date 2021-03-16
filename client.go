package socketcast

import "github.com/gorilla/websocket"

type Client struct {
	pool *Pool

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}
