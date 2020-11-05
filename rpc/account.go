package rpc

import "math/big"

// AccountBalance returns how many RAW is owned and how many have not yet been received by account.
func (c *Client) AccountBalance(account string) (balance, pending *big.Int, err error) {
	resp, err := c.send(map[string]string{"action": "account_balance", "account": account})
	if err != nil {
		return
	}
	if balance, err = toBig(resp["balance"]); err != nil {
		return
	}
	if pending, err = toBig(resp["pending"]); err != nil {
		return
	}
	return
}

// AccountBlockCount gets the number of blocks for a specific account.
func (c *Client) AccountBlockCount(account string) (blockCount uint64, err error) {
	resp, err := c.send(map[string]string{"action": "account_block_count", "account": account})
	if err != nil {
		return
	}
	if blockCount, err = toUint(resp["block_count"]); err != nil {
		return
	}
	return
}

// AccountGet gets the account number for the public key.
func (c *Client) AccountGet(key string) (account string, err error) {
	resp, err := c.send(map[string]string{"action": "account_get", "key": key})
	if err != nil {
		return
	}
	account = resp["account"]
	return
}

// AccountKey gets the public key for account.
func (c *Client) AccountKey(account string) (key string, err error) {
	resp, err := c.send(map[string]string{"action": "account_key", "account": account})
	if err != nil {
		return
	}
	key = resp["key"]
	return
}

// AccountRepresentative returns the representative for account.
func (c *Client) AccountRepresentative(account string) (representative string, err error) {
	resp, err := c.send(map[string]string{"action": "account_representative", "account": account})
	if err != nil {
		return
	}
	representative = resp["representative"]
	return
}

// AccountWeight returns the voting weight for account.
func (c *Client) AccountWeight(account string) (weight *big.Int, err error) {
	resp, err := c.send(map[string]string{"action": "account_weight", "account": account})
	if err != nil {
		return
	}
	if weight, err = toBig(resp["weight"]); err != nil {
		return
	}
	return
}
