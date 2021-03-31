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
	if config == nil {
		config = &Config{}
	}
	config.Defaultify()

	//set
	log := logger.NewLogger(config.LoggerLevel)

	//create hub
	pool := &Pool{
		hub: Hub{
			broadcast:  make(chan BroadcastMessage, config.Buffers.Broadcast),
			register:   make(chan *Client, config.Buffers.Register),
			unregister: make(chan *Client, config.Buffers.Unregister),
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
				if !message.Guard(client) {
					pool.Log.Infof("Skipped client %d", client.Conn.RemoteAddr().String())
					continue
				}
				if message.Msg == nil {
					if message.Generator == nil {
						pool.Log.Error("No messange and no generator for broadcast")
						continue
					}
					msg, should, err := message.Generator(client)
					if err != nil {
						pool.Log.Error("Error during generator in broadcast", err)
						continue
					}
					if !should {
						pool.Log.Debugf("skipped %s because generator ask for it", client.Conn.RemoteAddr().String())
						continue
					}
					message.Msg = msg
				}
				select {
				case client.send <- message.Msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (pool *Pool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := pool.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := newClient(pool, conn)
	pool.Log.Infof("Connected: %s", client.Conn.RemoteAddr().String())
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
		Msg: message, Guard: guard,
	}
}

func (pool *Pool) ForEach(generator func(*Client) (Message, bool, error), guard clientGuard) {
	if guard == nil {
		guard = func(c *Client) bool { return true }
	}
	pool.Log.Debug("About to run for each")
	pool.hub.broadcast <- BroadcastMessage{
		Guard: guard, Generator: func(c *Client) ([]byte, bool, error) {
			message, send, err := generator(c)
			if err != nil {
				pool.Log.Error(err)
				return nil, false, err
			}
			if !send {
				return nil, false, nil
			}
			d, err := json.Marshal(message)
			if err != nil {
				pool.Log.Error(err)
				return nil, false, err
			}
			return d, true, nil
		},
	}
}
