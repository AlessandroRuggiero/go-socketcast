package socketcast

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gobuffalo/logger"
	"github.com/gorilla/websocket"
)

type Pool struct {
	hub      Hub
	Log      logger.Logger
	Started  bool
	Config   Config
	Upgrader websocket.Upgrader
}

//CreatePool creates a new pool
func CreatePool(config *Config) *Pool {
	// prepare
	config.Defaultify()

	//set
	log := logger.NewLogger(config.LoggerLevel)

	//create hub
	pool := &Pool{
		hub: Hub{
			broadcast:  make(chan BroadcastMessage, config.HubBuffers.broadcast),
			register:   make(chan *Client, config.HubBuffers.register),
			unregister: make(chan *Client, config.HubBuffers.unregister),
			clients:    make(map[*Client]bool),
		},
		Log: log,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     config.CeckOrigin,
		},
	}
	pool.Config = *config
	//setup
	if !config.DisableAutostart {
		go pool.Start()
	} else {
		log.Info("You have left autostart off, remember to start the pool")
	}
	//return
	return pool
}

func (pool *Pool) Start() {
	if pool.Started {
		pool.Log.Error("pool already started")
		return
	}
	pool.Started = true
	pool.Log.Info("Pool started")
	pool.run()
}

func (pool *Pool) run() {
	h := &pool.hub
	for {
		select {
		case client := <-h.register:
			pool.Log.Debug("client registered")
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			pool.Log.Debug("message to broadcast")
			for client := range h.clients {
				if !message.guard(client) {
					pool.Log.Infof("Skipped client %d", client.Conn.RemoteAddr())
					continue
				}
				select {
				case client.send <- message.msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (pool *Pool) Serve(w http.ResponseWriter, r *http.Request) {
	conn, err := pool.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{Pool: pool, Conn: conn, send: make(chan []byte, 256), Metadata: make(map[string]interface{})}
	client.Start()
	pool.Log.Infof("Connesso: %s", client.Conn.RemoteAddr().String())
}

func (pool *Pool) Broadcast(msg interface{}, guard clientGuard) {
	if guard == nil {
		guard = func(c *Client) bool { return true }
	}
	message, err := json.Marshal(msg)
	if err != nil {
		pool.Log.Error(err)
		return
	}
	pool.Log.Debug("About to broadcast")
	pool.hub.broadcast <- BroadcastMessage{
		message, guard,
	}
}
