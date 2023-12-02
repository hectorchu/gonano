package rpc

import (
	"encoding/json"
	"time"
)

// AccountBalance returns how many RAW is owned and how many have not yet been received by account.
func (c *Client) AccountBalance(account string) (balance, receivable *RawAmount, err error) {
	resp, err := c.send(map[string]interface{}{"action": "account_balance", "account": account})
	if err != nil {
		return
	}
	var v struct{ Balance, Receivable *RawAmount }
	err = json.Unmarshal(resp, &v)
	return v.Balance, v.Receivable, err
}

// AccountBlockCount gets the number of blocks for a specific account.
func (c *Client) AccountBlockCount(account string) (blockCount uint64, err error) {
	resp, err := c.send(map[string]interface{}{"action": "account_block_count", "account": account})
	if err != nil {
		return
	}
	var v struct {
		BlockCount uint64 `json:"block_count,string"`
	}
	err = json.Unmarshal(resp, &v)
	return v.BlockCount, err
}

// AccountHistory reports send/receive information for an account.
func (c *Client) AccountHistory(account string, count int64, head BlockHash) (history []AccountHistory, previous BlockHash, err error) {
	body := map[string]interface{}{"action": "account_history", "account": account, "count": count}
	if head != nil {
		body["head"] = head
	}
	resp, err := c.send(body)
	if err != nil {
		return
	}
	var v struct {
		History  []AccountHistory
		Previous BlockHash
	}
	err = json.Unmarshal(resp, &v)
	return v.History, v.Previous, err
}

// AccountHistoryRaw reports all parameters of the block itself as seen in
// BlockCreate or other APIs returning blocks.
func (c *Client) AccountHistoryRaw(account string, count int64, head BlockHash) (history []AccountHistoryRaw, previous BlockHash, err error) {
	body := map[string]interface{}{"action": "account_history", "account": account, "count": count, "raw": true}
	if head != nil {
		body["head"] = head
	}
	resp, err := c.send(body)
	if err != nil {
		return
	}
	var v struct {
		History  []AccountHistoryRaw
		Previous BlockHash
	}
	err = json.Unmarshal(resp, &v)
	return v.History, v.Previous, err
}

// AccountInfo returns frontier, open block, change representative block,
// balance, last modified timestamp from local database & block count for
// account.
func (c *Client) AccountInfo(account string) (info AccountInfo, err error) {
	resp, err := c.send(map[string]interface{}{
		"action":         "account_info",
		"account":        account,
		"representative": true,
		"weight":         true,
		"receivable":     true,
	})
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &info)
	return
}

// AccountRepresentative returns the representative for account.
func (c *Client) AccountRepresentative(account string) (representative string, err error) {
	resp, err := c.send(map[string]interface{}{"action": "account_representative", "account": account})
	if err != nil {
		return
	}
	var v struct{ Representative string }
	err = json.Unmarshal(resp, &v)
	return v.Representative, err
}

// AccountWeight returns the voting weight for account.
func (c *Client) AccountWeight(account string) (weight *RawAmount, err error) {
	resp, err := c.send(map[string]interface{}{"action": "account_weight", "account": account})
	if err != nil {
		return
	}
	var v struct{ Weight *RawAmount }
	err = json.Unmarshal(resp, &v)
	return v.Weight, err
}

// AccountBalance returns how many RAW is owned and how many have not yet been received.
type AccountBalance struct {
	Balance, Receivable *RawAmount
}

// AccountsBalances returns how many RAW is owned and how many have not yet been received by accounts list.
func (c *Client) AccountsBalances(accounts []string) (balances map[string]*AccountBalance, err error) {
	resp, err := c.send(map[string]interface{}{"action": "accounts_balances", "accounts": accounts})
	if err != nil {
		return
	}
	var v struct{ Balances map[string]*AccountBalance }
	err = json.Unmarshal(resp, &v)
	return v.Balances, err
}

// AccountsFrontiers returns a list of pairs of account and block hash representing the head block for accounts list.
func (c *Client) AccountsFrontiers(accounts []string) (frontiers map[string]BlockHash, err error) {
	resp, err := c.send(map[string]interface{}{"action": "accounts_frontiers", "accounts": accounts})
	if err != nil {
		return
	}
	var u struct{ Frontiers string }
	if err = json.Unmarshal(resp, &u); err == nil && u.Frontiers == "" {
		return
	}
	var v struct{ Frontiers map[string]BlockHash }
	err = json.Unmarshal(resp, &v)
	return v.Frontiers, err
}

// AccountReceivable returns amount and source account.
type AccountReceivable struct {
	Amount *RawAmount
	Source string
}

// HashToReceivableMap maps receivable block hashes to amount and source account.
type HashToReceivableMap map[string]AccountReceivable

// UnmarshalJSON sets *h to a copy of data.
func (h *HashToReceivableMap) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err == nil && s == "" {
		return
	}
	var v map[string]AccountReceivable
	err = json.Unmarshal(data, &v)
	*h = v
	return
}

// AccountsReceivable returns a list of receivable block hashes with amount and source accounts.
func (c *Client) AccountsReceivable(accounts []string, count int64) (receivable map[string]HashToReceivableMap, err error) {
	resp, err := c.send(map[string]interface{}{
		"action":                 "accounts_receivable",
		"accounts":               accounts,
		"count":                  count,
		"include_only_confirmed": true,
		"source":                 true,
	})
	if err != nil {
		return
	}
	var u struct{ Blocks string }
	if err = json.Unmarshal(resp, &u); err == nil && u.Blocks == "" {
		return
	}
	var v struct {
		Blocks map[string]HashToReceivableMap
	}
	err = json.Unmarshal(resp, &v)
	return v.Blocks, err
}

// Delegators returns a list of pairs of delegator names given a representative account
// and its balance.
func (c *Client) Delegators(account string) (delegators map[string]*RawAmount, err error) {
	resp, err := c.send(map[string]interface{}{"action": "delegators", "account": account})
	if err != nil {
		return
	}
	var v struct{ Delegators map[string]*RawAmount }
	err = json.Unmarshal(resp, &v)
	return v.Delegators, err
}

// DelegatorsCount gets number of delegators for a specific representative account.
func (c *Client) DelegatorsCount(account string) (count uint64, err error) {
	resp, err := c.send(map[string]interface{}{"action": "delegators_count", "account": account})
	if err != nil {
		return
	}
	var v struct {
		Count uint64 `json:",string"`
	}
	err = json.Unmarshal(resp, &v)
	return v.Count, err
}

// FrontierCount reports the number of accounts in the ledger.
func (c *Client) FrontierCount() (count uint64, err error) {
	resp, err := c.send(map[string]interface{}{"action": "frontier_count"})
	if err != nil {
		return
	}
	var v struct {
		Count uint64 `json:",string"`
	}
	err = json.Unmarshal(resp, &v)
	return v.Count, err
}

// Frontiers returns a list of pairs of account and block hash representing the
// head block starting at account up to count.
func (c *Client) Frontiers(account string, count int64) (frontiers map[string]BlockHash, err error) {
	resp, err := c.send(map[string]interface{}{
		"action": "frontiers", "account": account, "count": count,
	})
	if err != nil {
		return
	}
	var v struct{ Frontiers map[string]BlockHash }
	err = json.Unmarshal(resp, &v)
	return v.Frontiers, err
}

// Ledger returns frontier, open block, change representative block, balance, last
// modified timestamp from local database & block count starting at account up to count.
func (c *Client) Ledger(account string, count int64, modifiedSince time.Time) (accounts map[string]AccountInfo, err error) {
	resp, err := c.send(map[string]interface{}{
		"action":         "ledger",
		"account":        account,
		"count":          count,
		"modified_since": modifiedSince.Unix(),
		"representative": true,
		"weight":         true,
		"receivable":     true,
	})
	if err != nil {
		return
	}
	var u struct{ Accounts string }
	if err = json.Unmarshal(resp, &u); err == nil && u.Accounts == "" {
		return
	}
	var v struct{ Accounts map[string]AccountInfo }
	err = json.Unmarshal(resp, &v)
	return v.Accounts, err
}

// Representatives returns a list of pairs of representative and its voting weight.
func (c *Client) Representatives(count int64) (representatives map[string]*RawAmount, err error) {
	resp, err := c.send(map[string]interface{}{"action": "representatives", "count": count})
	if err != nil {
		return
	}
	var v struct{ Representatives map[string]*RawAmount }
	err = json.Unmarshal(resp, &v)
	return v.Representatives, err
}

// Representative returns the weight of a representative.
type Representative struct{ Weight *RawAmount }

// RepresentativesOnline returns a list of online representative accounts that have voted recently.
func (c *Client) RepresentativesOnline() (representatives map[string]Representative, err error) {
	resp, err := c.send(map[string]interface{}{"action": "representatives_online", "weight": true})
	if err != nil {
		return
	}
	var v struct{ Representatives map[string]Representative }
	err = json.Unmarshal(resp, &v)
	return v.Representatives, err
}

// V23.0+ methods

// Receivable returns a list of block hashes which have not yet been received by this account.
func (c *Client) Receivable(
	account string, count int64, includeActive bool, threshold string,
) (receivable HashToReceivableMap, err error) {
	body := map[string]interface{}{
		"action":                 "receivable",
		"account":                account,
		"include_only_confirmed": true, // it defaults to false for v22.0 and below
		"source":                 true,
	}
	if count != 0 {
		body["count"] = count
	}
	// include_active defaults to false
	if includeActive != false {
		body["include_active"] = true
	}
	if threshold != "" {
		body["threshold"] = threshold
	}
	resp, err := c.send(body)
	if err != nil {
		return
	}
	var u struct{ Blocks string }
	if err = json.Unmarshal(resp, &u); err == nil && u.Blocks == "" {
		return
	}
	var v struct {
		Blocks HashToReceivableMap
	}
	err = json.Unmarshal(resp, &v)
	return v.Blocks, err
}
