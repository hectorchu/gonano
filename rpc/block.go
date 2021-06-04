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

// Chain returns a consecutive list of block hashes in the account chain starting
// at block back to count (direction from frontier back to open block, from newer
// blocks to older). Will list all blocks back to the open block of this chain when
// count is set to "-1". The requested block hash is included in the answer.
func (c *Client) Chain(block BlockHash, count int64) (blocks []BlockHash, err error) {
	resp, err := c.send(map[string]interface{}{"action": "chain", "block": block, "count": count})
	if err != nil {
		return
	}
	var v struct{ Blocks []BlockHash }
	err = json.Unmarshal(resp, &v)
	return v.Blocks, err
}

// Process publishes block to the network.
func (c *Client) Process(block *Block, subtype string) (hash BlockHash, err error) {
	resp, err := c.send(map[string]interface{}{
		"action":     "process",
		"json_block": true,
		"subtype":    subtype,
		"block":      block,
	})
	if err != nil {
		return
	}
	var v struct{ Hash BlockHash }
	err = json.Unmarshal(resp, &v)
	return v.Hash, err
}

// Republish rebroadcasts blocks starting at hash to the network.
func (c *Client) Republish(hash BlockHash, count, sources, destinations int64) (blocks []BlockHash, err error) {
	resp, err := c.send(map[string]interface{}{
		"action":       "republish",
		"hash":         hash,
		"count":        count,
		"sources":      sources,
		"destinations": destinations,
	})
	if err != nil {
		return
	}
	var v struct{ Blocks []BlockHash }
	err = json.Unmarshal(resp, &v)
	return v.Blocks, err
}

// Successors returns a consecutive list of block hashes in the account chain starting
// at block up to count (direction from open block up to frontier, from older
// blocks to newer). Will list all blocks up to frontier (latest block) of this chain
// when count is set to "-1". The requested block hash is included in the answer.
func (c *Client) Successors(block BlockHash, count int64) (blocks []BlockHash, err error) {
	resp, err := c.send(map[string]interface{}{"action": "successors", "block": block, "count": count})
	if err != nil {
		return
	}
	var v struct{ Blocks []BlockHash }
	err = json.Unmarshal(resp, &v)
	return v.Blocks, err
}
