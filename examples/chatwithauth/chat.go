package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlessandroRuggiero/go-socketcast"
)

func checkAuth(data string) (bool, uint64, string, error) {
	return true, 42, "example", nil // should check the auth (like on the db)
}
func guard(c *socketcast.Client) bool { // the guard is used to exclude clients from reciving the message
	return c.Auth.IsAuthenticated() // the message witll be ste only to clients that are authenticated
}

func handle(c *socketcast.Client, msg []byte) bool {
	var message socketcast.Message       // Any json tagged struct is ok, message is just a redy to use struct provided by the library
	err := json.Unmarshal(msg, &message) // parse the message
	if err != nil {
		c.Pool.Log.Warnf("Error parsing data %v", msg)
		return true //close the connection on malformed message
	}
	switch message.Type { // suggested not to use 0
	case 1: // code for message (complitly arbitrary)
		if !c.Auth.IsAuthenticated() {
			return true // terminate connectiomn if not authenticated user sends a message
		}
		c.Pool.Broadcast(message, guard) // use a guard to send data only to authenticated users
	case 2: // code for authenticate (complitly arbitrary)
		ok, id, token, err := checkAuth(message.Body)
		if err != nil {
			c.Send(socketcast.Message{
				Type: 3, // code for server error
				Body: "Impossible to authenticate, server error",
			})
			return false
		}
		if !ok {
			c.Pool.Log.Info("Auth failed: ", c.Conn.RemoteAddr().String())
			return true //terminate connection
		}
		c.Auth.Authenticate() // update the client auth
		c.Auth.SetToken(token)
		c.Auth.SetId(id)
	}
	return false // true if you want to close the connection
}

func main() {
	pool := socketcast.CreatePool(&socketcast.Config{
		OnMessage: handle, // callback to call when a message is recived
		OnConnect: func(c *socketcast.Client) { // callback to call when a client connects
			log.Println(c.Conn.RemoteAddr().String(), "entered the chat")
		},
		OnDisconnect: func(c *socketcast.Client) { // callback to call when a client disconnects
			log.Println(c.Conn.RemoteAddr().String(), "left the chat")
		},
	})
	http.Handle("/chat", pool) //expose socket connection endpoint
	log.Fatal(http.ListenAndServe(":8080", nil))
}
