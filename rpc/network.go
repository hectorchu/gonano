package rpc

import "encoding/json"

// AvailableSupply returns how many raw are in the public supply.
func (c *Client) AvailableSupply() (available *RawAmount, err error) {
	resp, err := c.send(map[string]interface{}{"action": "available_supply"})
	if err != nil {
		return
	}
	var v struct{ Available *RawAmount }
	err = json.Unmarshal(resp, &v)
	return v.Available, err
}
