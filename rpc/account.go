package rpc

import (
	"errors"
	"math/big"
)

// AccountBalance returns how many RAW is owned and how many have not yet been received by account.
func (c *Client) AccountBalance(account string) (balance, pending *big.Int, err error) {
	resp, err := c.send(map[string]interface{}{"action": "account_balance", "account": account})
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
	resp, err := c.send(map[string]interface{}{"action": "account_block_count", "account": account})
	if err != nil {
		return
	}
	return toUint(resp["block_count"])
}

// AccountGet gets the account number for the public key.
func (c *Client) AccountGet(key string) (account string, err error) {
	resp, err := c.send(map[string]interface{}{"action": "account_get", "key": key})
	if err != nil {
		return
	}
	return toStr(resp["account"])
}

// AccountHistory reports send/receive information for an account.
func (c *Client) AccountHistory(account string, count int64, head string) (history []AccountHistory, previous []byte, err error) {
	body := map[string]interface{}{"action": "account_history", "account": account, "count": count}
	if head != "" {
		body["head"] = head
	}
	resp, err := c.send(body)
	if err != nil {
		return
	}
	h, ok := resp["history"].([]interface{})
	if !ok {
		err = errors.New("failed to cast history array")
		return
	}
	history = make([]AccountHistory, len(h))
	for i, h := range h {
		h, ok := h.(map[string]interface{})
		if !ok {
			err = errors.New("failed to cast history array")
			return
		}
		if err = history[i].parse(h); err != nil {
			return
		}
	}
	if v, ok := resp["previous"]; ok {
		if previous, err = toBytes(v); err != nil {
			return
		}
	}
	return
}

// AccountHistoryRaw reports all parameters of the block itself as seen in
// BlockCreate or other APIs returning blocks.
func (c *Client) AccountHistoryRaw(account string, count int64, head string) (history []AccountHistoryRaw, previous []byte, err error) {
	body := map[string]interface{}{"action": "account_history", "account": account, "count": count, "raw": true}
	if head != "" {
		body["head"] = head
	}
	resp, err := c.send(body)
	if err != nil {
		return
	}
	h, ok := resp["history"].([]interface{})
	if !ok {
		err = errors.New("failed to cast history array")
		return
	}
	history = make([]AccountHistoryRaw, len(h))
	for i, h := range h {
		h, ok := h.(map[string]interface{})
		if !ok {
			err = errors.New("failed to cast history array")
			return
		}
		if err = history[i].parse(h); err != nil {
			return
		}
	}
	if v, ok := resp["previous"]; ok {
		if previous, err = toBytes(v); err != nil {
			return
		}
	}
	return
}

// AccountInfo returns frontier, open block, change representative block,
// balance, last modified timestamp from local database & block count for
// account.
func (c *Client) AccountInfo(account string) (info AccountInfo, err error) {
	resp, err := c.send(map[string]interface{}{"action": "account_info", "account": account})
	if err != nil {
		return
	}
	err = info.parse(resp)
	return
}

// AccountKey gets the public key for account.
func (c *Client) AccountKey(account string) (key []byte, err error) {
	resp, err := c.send(map[string]interface{}{"action": "account_key", "account": account})
	if err != nil {
		return
	}
	return toBytes(resp["key"])
}

// AccountRepresentative returns the representative for account.
func (c *Client) AccountRepresentative(account string) (representative string, err error) {
	resp, err := c.send(map[string]interface{}{"action": "account_representative", "account": account})
	if err != nil {
		return
	}
	return toStr(resp["representative"])
}

// AccountWeight returns the voting weight for account.
func (c *Client) AccountWeight(account string) (weight *big.Int, err error) {
	resp, err := c.send(map[string]interface{}{"action": "account_weight", "account": account})
	if err != nil {
		return
	}
	return toBig(resp["weight"])
}

// AccountBalance returns how many RAW is owned and how many have not yet been received.
type AccountBalance struct {
	Balance, Pending *big.Int
}

// AccountsBalances returns how many RAW is owned and how many have not yet been received by accounts list.
func (c *Client) AccountsBalances(accounts []string) (balances map[string]AccountBalance, err error) {
	resp, err := c.send(map[string]interface{}{"action": "accounts_balances", "accounts": accounts})
	if err != nil {
		return
	}
	b, ok := resp["balances"].(map[string]interface{})
	if !ok {
		err = errors.New("failed to cast balances map")
		return
	}
	balances = make(map[string]AccountBalance)
	for account, b := range b {
		b, ok := b.(map[string]interface{})
		if !ok {
			err = errors.New("failed to cast balances map")
			return
		}
		var balance AccountBalance
		if balance.Balance, err = toBig(b["balance"]); err != nil {
			return
		}
		if balance.Pending, err = toBig(b["pending"]); err != nil {
			return
		}
		balances[account] = balance
	}
	return
}
