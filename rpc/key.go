package rpc

import "encoding/json"

// DeterministicKey derives a deterministic keypair from seed based on index.
func (c *Client) DeterministicKey(seed HexData, index uint64) (
	private, public HexData, account string, err error,
) {
	resp, err := c.send(map[string]interface{}{
		"action": "deterministic_key", "seed": seed, "index": index,
	})
	if err != nil {
		return
	}
	var v struct {
		Private, Public HexData
		Account         string
	}
	err = json.Unmarshal(resp, &v)
	return v.Private, v.Public, v.Account, err
}

// KeyCreate generates an adhoc random keypair.
func (c *Client) KeyCreate() (private, public HexData, account string, err error) {
	resp, err := c.send(map[string]interface{}{"action": "key_create"})
	if err != nil {
		return
	}
	var v struct {
		Private, Public HexData
		Account         string
	}
	err = json.Unmarshal(resp, &v)
	return v.Private, v.Public, v.Account, err
}

// KeyExpand derives public key and account number from private key.
func (c *Client) KeyExpand(key HexData) (private, public HexData, account string, err error) {
	resp, err := c.send(map[string]interface{}{"action": "key_expand", "key": key})
	if err != nil {
		return
	}
	var v struct {
		Private, Public HexData
		Account         string
	}
	err = json.Unmarshal(resp, &v)
	return v.Private, v.Public, v.Account, err
}
