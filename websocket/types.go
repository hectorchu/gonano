package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hectorchu/gonano/rpc"
)

// Confirmation reports a block confirmation.
type Confirmation struct {
	Time    time.Time
	Account string
	Amount  *rpc.RawAmount
	Hash    rpc.BlockHash
	Type    string `json:"confirmation_type"`
	Block   *rpc.Block
}

type message struct{ m interface{} }

// UnmarshalJSON sets *m to a copy of data.
func (m *message) UnmarshalJSON(data []byte) (err error) {
	var v struct {
		Topic string
		Time  int64 `json:",string"`
	}
	if err = json.Unmarshal(data, &v); err != nil {
		return
	}
	switch v.Topic {
	case "confirmation":
		var w struct{ Message *Confirmation }
		if err = json.Unmarshal(data, &w); err != nil {
			return
		}
		w.Message.Time = time.Unix(0, v.Time*1e6).UTC()
		*m = message{m: w.Message}
	default:
		err = errors.New(fmt.Sprint("unknown topic", v.Topic))
	}
	return
}
