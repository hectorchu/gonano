package rpc

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/hectorchu/gonano/util"
	"golang.org/x/crypto/blake2b"
)

// AccountHistory reports send/receive information within a block.
type AccountHistory struct {
	Type           string     `json:"type"`
	Account        string     `json:"account"`
	Amount         *RawAmount `json:"amount"`
	LocalTimestamp uint64     `json:"local_timestamp,string"`
	Height         uint64     `json:"height,string"`
	Hash           BlockHash  `json:"hash"`
}

// AccountHistoryRaw reports all parameters of the block itself as seen in
// BlockCreate or other APIs returning blocks.
type AccountHistoryRaw struct {
	Type           string     `json:"type"`
	Representative string     `json:"representative"`
	Link           BlockHash  `json:"link"`
	Balance        *RawAmount `json:"balance"`
	Previous       BlockHash  `json:"previous"`
	Subtype        string     `json:"subtype"`
	Account        string     `json:"account"`
	Amount         *RawAmount `json:"amount"`
	LocalTimestamp uint64     `json:"local_timestamp,string"`
	Height         uint64     `json:"height,string"`
	Hash           BlockHash  `json:"hash"`
	Work           HexData    `json:"work"`
	Signature      HexData    `json:"signature"`
}

// AccountInfo returns frontier, open block, change representative block,
// balance, last modified timestamp from local database & block count for
// account.
type AccountInfo struct {
	Frontier                   BlockHash  `json:"frontier"`
	OpenBlock                  BlockHash  `json:"open_block"`
	RepresentativeBlock        BlockHash  `json:"representative_block"`
	Balance                    *RawAmount `json:"balance"`
	ModifiedTimestamp          uint64     `json:"modified_timestamp,string"`
	BlockCount                 uint64     `json:"block_count,string"`
	ConfirmationHeight         uint64     `json:"confirmation_height,string"`
	ConfirmationHeightFrontier BlockHash  `json:"confirmation_height_frontier"`
	AccountVersion             uint64     `json:"account_version,string"`
	Representative             string     `json:"representative"`
	Weight                     *RawAmount `json:"weight"`
	Receivable                 *RawAmount `json:"receivable"`
}

// Block corresponds to the JSON representation of a block.
type Block struct {
	Type           string     `json:"type"`
	Account        string     `json:"account"`
	Previous       BlockHash  `json:"previous"`
	Representative string     `json:"representative"`
	Balance        *RawAmount `json:"balance"`
	Link           BlockHash  `json:"link"`
	LinkAsAccount  string     `json:"link_as_account"`
	Signature      HexData    `json:"signature"`
	Work           HexData    `json:"work"`
}

// Hash calculates the block hash.
func (b *Block) Hash() (hash BlockHash, err error) {
	h, err := blake2b.New256(nil)
	if err != nil {
		return
	}
	h.Write(make([]byte, 31))
	h.Write([]byte{6})
	pubkey, err := util.AddressToPubkey(b.Account)
	if err != nil {
		return
	}
	h.Write(pubkey)
	h.Write(b.Previous)
	pubkey, err = util.AddressToPubkey(b.Representative)
	if err != nil {
		return
	}
	h.Write(pubkey)
	h.Write(b.Balance.FillBytes(make([]byte, 16)))
	h.Write(b.Link)
	return h.Sum(nil), nil
}

// BlockHash represents a block hash.
type BlockHash []byte

func (h BlockHash) String() string {
	return strings.ToUpper(hex.EncodeToString(h))
}

// MarshalJSON returns the JSON encoding of h.
func (h BlockHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON sets *h to a copy of data.
func (h *BlockHash) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	*h, err = hex.DecodeString(s)
	return
}

// BlockInfo retrieves a json representation of a block.
type BlockInfo struct {
	BlockAccount   string     `json:"block_account"`
	Amount         *RawAmount `json:"amount"`
	Balance        *RawAmount `json:"balance"`
	Height         uint64     `json:"height,string"`
	LocalTimestamp uint64     `json:"local_timestamp,string"`
	Confirmed      bool       `json:"confirmed,string"`
	Contents       *Block     `json:"contents"`
	Subtype        string     `json:"subtype"`
}

// HexData represents generic hex data.
type HexData []byte

// MarshalJSON returns the JSON encoding of h.
func (h HexData) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(h))
}

// UnmarshalJSON sets *h to a copy of data.
func (h *HexData) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	*h, err = hex.DecodeString(s)
	return
}

// RawAmount represents an amount of nano in RAWs.
type RawAmount struct{ big.Int }

// MarshalJSON returns the JSON encoding of r.
func (r *RawAmount) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON sets *r to a copy of data.
func (r *RawAmount) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	if _, ok := r.SetString(s, 10); !ok {
		err = errors.New("unable to parse amount")
	}
	return
}
