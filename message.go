package socketcast

import "encoding/json"

type Message struct {
	Type     int                    `json:"type"`
	Body     string                 `json:"body"`
	Metadata map[string]interface{} `json:"-"`
}

func (msg Message) Read() (message []byte, err error) {
	message, err = json.Marshal(msg)
	return
}
func (msg Message) GetType() int {
	return msg.Type
}

func (msg Message) GetMetadata() map[string]interface{} {
	return msg.Metadata
}
