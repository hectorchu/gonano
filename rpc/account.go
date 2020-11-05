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
