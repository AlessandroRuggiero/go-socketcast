package main

import (
	"log"
	"net/http"

	"github.com/AlessandroRuggiero/go-socketcast"
)

func echo(c *socketcast.Client, msg []byte) bool {
	c.Send(socketcast.Message{
		Type: 1, // Complitly arbitrary, just for better clarity with client
		Body: string(msg),
	})
	return false // true if you want to close the connection
}

func main() {
	pool := socketcast.CreatePool(&socketcast.Config{
		OnMessage: echo, // callback to call when a message is recived
	})
	http.Handle("/echo", pool)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
