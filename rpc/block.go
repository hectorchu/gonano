package rpc

import (
	"encoding/json"
)

// BlockAccount returns the account containing block.
func (c *Client) BlockAccount(hash BlockHash) (account string, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_account", "hash": hash})
	if err != nil {
		return
	}
	var v struct{ Account string }
	err = json.Unmarshal(resp, &v)
	return v.Account, err
}

// BlockConfirm requests confirmation for block from known online representative nodes.
func (c *Client) BlockConfirm(hash BlockHash) (started bool, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_confirm", "hash": hash})
	if err != nil {
		return
	}
	var v struct {
		Started int `json:",string"`
	}
	err = json.Unmarshal(resp, &v)
	return v.Started == 1, err
}

// BlockCount reports the number of blocks in the ledger and unchecked synchronizing blocks.
func (c *Client) BlockCount() (cemented, count, unchecked uint64, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_count"})
	if err != nil {
		return
	}
	var v struct {
		Cemented, Count, Unchecked uint64 `json:",string"`
	}
	err = json.Unmarshal(resp, &v)
	return v.Cemented, v.Count, v.Unchecked, err
}

// BlockCountType reports the number of blocks in the ledger by type.
func (c *Client) BlockCountType() (send, receive, open, change, state uint64, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_count_type"})
	if err != nil {
		return
	}
	var v struct {
		Send, Receive, Open, Change, State uint64 `json:",string"`
	}
	err = json.Unmarshal(resp, &v)
	return v.Send, v.Receive, v.Open, v.Change, v.State, err
}

// BlockHash returns the block hash for the given block content.
func (c *Client) BlockHash(block *Block) (hash BlockHash, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_hash", "json_block": true, "block": block})
	if err != nil {
		return
	}
	var v struct{ Hash BlockHash }
	err = json.Unmarshal(resp, &v)
	return v.Hash, err
}

// BlockInfo retrieves a json representation of a block.
func (c *Client) BlockInfo(hash BlockHash) (info BlockInfo, err error) {
	resp, err := c.send(map[string]interface{}{"action": "block_info", "json_block": true, "hash": hash})
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &info)
	return
}

// Blocks retrieves a json representations of blocks.
func (c *Client) Blocks(hashes []BlockHash) (blocks map[string]*Block, err error) {
	resp, err := c.send(map[string]interface{}{"action": "blocks", "json_block": true, "hashes": hashes})
	if err != nil {
		return
	}
	var v struct{ Blocks map[string]*Block }
	err = json.Unmarshal(resp, &v)
	return v.Blocks, err
}

// BlocksInfo retrieves a json representations of blocks in contents.
func (c *Client) BlocksInfo(hashes []BlockHash) (blocks map[string]*BlockInfo, err error) {
	resp, err := c.send(map[string]interface{}{"action": "blocks_info", "json_block": true, "hashes": hashes})
	if err != nil {
		return
	}
	var v struct{ Blocks map[string]*BlockInfo }
	err = json.Unmarshal(resp, &v)
	return v.Blocks, err
}
