package socketcast

import "net/http"

//Config holds config data for a new pool
type Config struct {
	LoggerLevel            string
	HubBuffers             BufferSize
	DisableAutostart       bool
	DisableClientAutostart bool
	OnMessage              func(c *Client, msg []byte) bool
	CeckOrigin             func(r *http.Request) bool
}

func (config *Config) Defaultify() {
	if len(config.LoggerLevel) == 0 {
		config.LoggerLevel = "debug"
	}
	if config.OnMessage == nil {
		config.OnMessage = func(c *Client, msg []byte) bool {
			c.Pool.Log.Debug(string(msg))
			return false
		}
	}
	if config.CeckOrigin == nil {
		config.CeckOrigin = func(r *http.Request) bool { return true }
	}
}
