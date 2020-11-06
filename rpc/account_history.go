package rpc

import (
	"math/big"
	"time"
)

// AccountHistory reports send/receive information within a block.
type AccountHistory struct {
	Type           string
	Account        string
	Amount         *big.Int
	LocalTimestamp time.Time
	Height         uint64
	Hash           []byte
}

func (h *AccountHistory) parse(x map[string]interface{}) (err error) {
	if h.Type, err = toStr(x["type"]); err != nil {
		return
	}
	if h.Account, err = toStr(x["account"]); err != nil {
		return
	}
	if h.Amount, err = toBig(x["amount"]); err != nil {
		return
	}
	if h.LocalTimestamp, err = toTime(x["local_timestamp"]); err != nil {
		return
	}
	if h.Height, err = toUint(x["height"]); err != nil {
		return
	}
	if h.Hash, err = toBytes(x["hash"]); err != nil {
		return
	}
	return
}

// AccountHistoryRaw reports all parameters of the block itself as seen in
// BlockCreate or other APIs returning blocks.
type AccountHistoryRaw struct {
	Type           string
	Representative string
	Link           []byte
	Balance        *big.Int
	Previous       []byte
	Subtype        string
	Account        string
	Amount         *big.Int
	LocalTimestamp time.Time
	Height         uint64
	Hash           []byte
	Work           []byte
	Signature      []byte
}

func (h *AccountHistoryRaw) parse(x map[string]interface{}) (err error) {
	if h.Type, err = toStr(x["type"]); err != nil {
		return
	}
	if h.Representative, err = toStr(x["representative"]); err != nil {
		return
	}
	if h.Link, err = toBytes(x["link"]); err != nil {
		return
	}
	if h.Balance, err = toBig(x["balance"]); err != nil {
		return
	}
	if h.Previous, err = toBytes(x["previous"]); err != nil {
		return
	}
	if h.Subtype, err = toStr(x["subtype"]); err != nil {
		return
	}
	if h.Account, err = toStr(x["account"]); err != nil {
		return
	}
	if h.Amount, err = toBig(x["amount"]); err != nil {
		return
	}
	if h.LocalTimestamp, err = toTime(x["local_timestamp"]); err != nil {
		return
	}
	if h.Height, err = toUint(x["height"]); err != nil {
		return
	}
	if h.Hash, err = toBytes(x["hash"]); err != nil {
		return
	}
	if h.Work, err = toBytes(x["work"]); err != nil {
		return
	}
	if h.Signature, err = toBytes(x["signature"]); err != nil {
		return
	}
	return
}
