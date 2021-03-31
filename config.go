package socketcast

import (
	"net/http"
)

//Config holds config data for a new pool
type Config struct {
	LoggerLevel            string
	Buffers                BufferSize
	DisableAutostart       bool
	DisableClientAutostart bool
	OnMessage              func(c *Client, msg []byte) bool
	OnConnect              func(c *Client)
	OnDisconnect           func(c *Client)
	CeckOrigin             func(r *http.Request) bool
}

func (config *Config) Defaultify() {
	if len(config.LoggerLevel) == 0 {
		config.LoggerLevel = "debug"
	}
	if config.OnMessage == nil {
		config.OnMessage = defOnMessage
	}
	if config.CeckOrigin == nil {
		config.CeckOrigin = defCeckOrigin
	}
	if config.OnConnect == nil {
		config.OnConnect = defOnConnect
	}
	if config.OnDisconnect == nil {
		config.OnDisconnect = defOnDisconnect
	}
	if config.Buffers.Send == 0 {
		config.Buffers.Send = 256
	}
}
