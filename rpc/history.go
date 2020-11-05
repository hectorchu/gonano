package rpc

import "math/big"

// History reports send/receive information within a block.
type History struct {
	Type           string
	Account        string
	Amount         *big.Int
	LocalTimestamp uint64
	Height         uint64
	Hash           string
}

func (h *History) parse(x map[string]interface{}) (err error) {
	if h.Type, err = toStr(x["type"]); err != nil {
		return
	}
	if h.Account, err = toStr(x["account"]); err != nil {
		return
	}
	if h.Amount, err = toBig(x["amount"]); err != nil {
		return
	}
	if h.LocalTimestamp, err = toUint(x["local_timestamp"]); err != nil {
		return
	}
	if h.Height, err = toUint(x["height"]); err != nil {
		return
	}
	if h.Hash, err = toStr(x["hash"]); err != nil {
		return
	}
	return
}
