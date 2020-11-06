package rpc

import "encoding/hex"

// BlockAccount returns the account containing block.
func (c *Client) BlockAccount(hash []byte) (account string, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_account", "hash": hex.EncodeToString(hash)})
	if err != nil {
		return
	}
	return toStr(resp["account"])
}

// BlockConfirm requests confirmation for block from known online representative nodes.
func (c *Client) BlockConfirm(hash []byte) (started bool, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_confirm", "hash": hex.EncodeToString(hash)})
	if err != nil {
		return
	}
	var started2 uint64
	if started2, err = toUint(resp["started"]); err != nil {
		return
	}
	return started2 == 1, nil
}

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

// BlockCountType reports the number of blocks in the ledger by type
// (send, receive, open, change, state with version).
func (c *Client) BlockCountType() (send, receive, open, change, state uint64, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_count_type"})
	if err != nil {
		return
	}
	if send, err = toUint(resp["send"]); err != nil {
		return
	}
	if receive, err = toUint(resp["receive"]); err != nil {
		return
	}
	if open, err = toUint(resp["open"]); err != nil {
		return
	}
	if change, err = toUint(resp["change"]); err != nil {
		return
	}
	if state, err = toUint(resp["state"]); err != nil {
		return
	}
	return
}

// BlockHash returns the block hash for the given block content.
func (c *Client) BlockHash(block *Block) (hash []byte, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_hash", "json_block": true, "block": block})
	if err != nil {
		return
	}
	return toBytes(resp["hash"])
}
