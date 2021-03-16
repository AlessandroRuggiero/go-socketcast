package main

import (
	"log"
	"net/http"

	"github.com/AlessandroRuggiero/go-socketcast"
)

func onm(c *socketcast.Client, msg []byte) bool {
	c.Pool.Log.Debug(string(msg))
	c.Pool.Broadcast(socketcast.Message{
		Body: string(msg),
	}, nil)
	return false
}

func main() {
	pool := socketcast.CreatePool(&socketcast.Config{
		OnMessage: onm,
	})
	http.HandleFunc("/ws", pool.Serve)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
