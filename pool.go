package socketcast

import (
	"github.com/gobuffalo/logger"
)

type Pool struct {
	hub     Hub
	log     logger.Logger
	Started bool
}

//CreatePool creates a new pool
func CreatePool(config *Config) *Pool {
	// prepare
	config.Defoultify()

	//set
	log := logger.NewLogger(config.loggerLevel)

	//create hub
	pool := &Pool{
		hub: Hub{
			broadcast:  make(chan []byte, config.hubBuffers.broadcast),
			register:   make(chan *Client, config.hubBuffers.register),
			unregister: make(chan *Client, config.hubBuffers.unregister),
			clients:    make(map[*Client]bool),
		},
		log: log,
	}
	//setup
	if !config.disableAutostart {
		go pool.Start()
	} else {
		log.Info("You have left autostart off, remember to start the pool")
	}
	//return
	return pool
}

func (pool *Pool) Start() {
	if pool.Started {
		pool.log.Error("pool already started")
		return
	}
	pool.Started = true
	pool.log.Info("Pool started")
}
