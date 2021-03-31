package socketcast

import "github.com/gorilla/websocket"

func analizeSocketError(c *Client, err *websocket.CloseError) {
	switch err.Code {
	case websocket.CloseAbnormalClosure:
		c.Pool.Log.Infof("Client %s disconnected (AbnormalClosure)", c.Conn.RemoteAddr().String())
	case websocket.CloseGoingAway, websocket.CloseNormalClosure:
		c.Pool.Log.Debug("Detected close ")
		// nothing to do, just a normal close
	case websocket.CloseNoStatusReceived:
		c.Pool.Log.Infof("Client %s disconnected (no status)", c.Conn.RemoteAddr().String())
	default:
		c.Pool.Log.Warnf("UnexpectedCloseError, closing connection: %s", err.Error())
	}
}
