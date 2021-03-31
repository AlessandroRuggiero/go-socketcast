package main

import (
	"log"
	"net/http"

	"github.com/AlessandroRuggiero/go-socketcast"
)

func main() {
	pool := socketcast.CreatePool(nil)
	http.Handle("/ws", pool)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
