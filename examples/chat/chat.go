package main

import (
	"log"
	"net/http"

	"github.com/AlessandroRuggiero/go-socketcast"
)

func handleMessage(c *socketcast.Client, msg []byte) bool {
	c.Pool.Broadcast( // broadcast the message
		socketcast.Message{ // Any json tagged struct is ok, message is just a redy to use struct provided by the library
			Type: 2,           // Complitly arbitrary, just for better clarity with client
			Body: string(msg), // sent the message
		}, nil)
	return false // true if you want to close the connection
}

func main() {
	pool := socketcast.CreatePool(&socketcast.Config{
		OnMessage: handleMessage,
	})
	http.Handle("/ws", pool)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
