package rpc

import "encoding/hex"

// BlockCount reports the number of blocks in the ledger and unchecked synchronizing blocks.
func (c *Client) BlockCount() (cemented, count, unchecked uint64, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_count"})
	if err != nil {
		return
	}
	if cemented, err = toUint(resp["cemented"]); err != nil {
		return
	}
	if count, err = toUint(resp["count"]); err != nil {
		return
	}
	if unchecked, err = toUint(resp["unchecked"]); err != nil {
		return
	}
	return
}

// BlockAccount returns the account containing block.
func (c *Client) BlockAccount(hash []byte) (account string, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_account", "hash": hex.EncodeToString(hash)})
	if err != nil {
		return
	}
	return toStr(resp["account"])
}