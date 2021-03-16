package socketcast

//Config holds config data for a new pool
type Config struct {
	LoggerLevel      string
	HubBuffers       HubBuffer
	DisableAutostart bool
	OnMessage        func(c *Client, msg []byte) bool
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
}
