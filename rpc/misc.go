package rpc

import (
	"encoding/json"
	"strconv"
)

// ActiveDifficulty returns the difficulty values (16 hexadecimal digits string, 64 bit)
// for the minimum required on the network (network_minimum) as well as the current active
// difficulty seen on the network (network_current, 10 second trended average of adjusted
// difficulty seen on prioritized transactions) which can be used to perform rework for
// better prioritization of transaction processing. A multiplier of the network_current
// from the base difficulty of network_minimum is also provided for comparison.
// network_receive_minimum and network_receive_current are also provided as lower
// thresholds exclusively for receive blocks.
func (c *Client) ActiveDifficulty() (
	multiplier float64,
	networkCurrent, networkMinimum,
	networkReceiveCurrent, networkReceiveMinimum HexData,
	difficultyTrend []float64,
	err error,
) {
	resp, err := c.send(map[string]interface{}{"action": "active_difficulty", "include_trend": true})
	if err != nil {
		return
	}
	var v struct {
		Multiplier            float64  `json:"multiplier,string"`
		NetworkCurrent        HexData  `json:"network_current"`
		NetworkMinimum        HexData  `json:"network_minimum"`
		NetworkReceiveCurrent HexData  `json:"network_receive_current"`
		NetworkReceiveMinimum HexData  `json:"network_receive_minimum"`
		DifficultyTrend       []string `json:"difficulty_trend"`
	}
	if err = json.Unmarshal(resp, &v); err != nil {
		return
	}
	difficultyTrend = make([]float64, len(v.DifficultyTrend))
	for i, s := range v.DifficultyTrend {
		if difficultyTrend[i], err = strconv.ParseFloat(s, 64); err != nil {
			return
		}
	}
	return v.Multiplier,
		v.NetworkCurrent, v.NetworkMinimum,
		v.NetworkReceiveCurrent, v.NetworkReceiveMinimum,
		difficultyTrend, err
}

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
