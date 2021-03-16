package socketcast

//Config holds config data for a new pool
type Config struct {
	loggerLevel      string
	hubBuffers       HubBuffer
	disableAutostart bool
}

func (config *Config) Defoultify() {
	if len(config.loggerLevel) == 0 {
		config.loggerLevel = "debug"
	}
}
