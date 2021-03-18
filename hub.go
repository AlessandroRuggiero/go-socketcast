package socketcast

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
}

type BroadcastMessage struct {
	Msg       []byte
	Guard     clientGuard
	Generator func(*Client) ([]byte, error)
}
