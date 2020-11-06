package rpc

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
)

// Block corresponds to the JSON representation of a block.
type Block struct {
	Type           string
	Account        string
	Previous       []byte
	Representative string
	Balance        *big.Int
	Link           []byte
	LinkAsAccount  string
	Signature      []byte
	Work           []byte
}

func (b *Block) parse(x map[string]interface{}) (err error) {
	if b.Type, err = toStr(x["type"]); err != nil {
		return
	}
	if b.Account, err = toStr(x["account"]); err != nil {
		return
	}
	if b.Previous, err = toBytes(x["previous"]); err != nil {
		return
	}
	if b.Representative, err = toStr(x["representative"]); err != nil {
		return
	}
	if b.Balance, err = toBig(x["balance"]); err != nil {
		return
	}
	if b.Link, err = toBytes(x["link"]); err != nil {
		return
	}
	if b.LinkAsAccount, err = toStr(x["link_as_account"]); err != nil {
		return
	}
	if b.Signature, err = toBytes(x["signature"]); err != nil {
		return
	}
	if b.Work, err = toBytes(x["work"]); err != nil {
		return
	}
	return
}

// MarshalJSON returns the JSON encoding of b.
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":            b.Type,
		"account":         b.Account,
		"previous":        hex.EncodeToString(b.Previous),
		"representative":  b.Representative,
		"balance":         b.Balance.String(),
		"link":            hex.EncodeToString(b.Link),
		"link_as_account": b.LinkAsAccount,
		"signature":       hex.EncodeToString(b.Signature),
		"work":            hex.EncodeToString(b.Work),
	})
}
