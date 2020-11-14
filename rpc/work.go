package rpc

import "encoding/json"

// WorkCancel stops generating work for block.
func (c *Client) WorkCancel(hash BlockHash) (err error) {
	_, err = c.send(map[string]interface{}{"action": "work_cancel", "hash": hash})
	return
}

// WorkGenerate generates work for block. hash is the frontier of the account
// or in the case of an open block, the public key representation of the account.
func (c *Client) WorkGenerate(hash BlockHash, difficulty HexData) (
	work, difficulty2 HexData, multiplier float64, err error,
) {
	resp, err := c.send(map[string]interface{}{
		"action": "work_generate", "hash": hash, "difficulty": difficulty,
	})
	if err != nil {
		return
	}
	var v struct {
		Work, Difficulty HexData
		Multiplier       float64 `json:",string"`
	}
	err = json.Unmarshal(resp, &v)
	return v.Work, v.Difficulty, v.Multiplier, err
}

// WorkValidate checks whether work is valid for block. Provides two values:
// validAll is true if the work is valid at the current network difficulty
// (work can be used for any block).
// validReceive is true if the work is valid for use in a receive block.
func (c *Client) WorkValidate(hash BlockHash, work HexData) (
	validAll, validReceive bool,
	difficulty HexData, multiplier float64, err error,
) {
	resp, err := c.send(map[string]interface{}{"action": "work_validate", "hash": hash, "work": work})
	if err != nil {
		return
	}
	var v struct {
		ValidAll     int `json:"valid_all,string"`
		ValidReceive int `json:"valid_receive,string"`
		Difficulty   HexData
		Multiplier   float64 `json:",string"`
	}
	err = json.Unmarshal(resp, &v)
	return v.ValidAll == 1, v.ValidReceive == 1, v.Difficulty, v.Multiplier, err
}
