package wallet

import (
	"encoding/hex"

	"github.com/hectorchu/gonano/rpc"
	"github.com/hectorchu/gonano/wallet/ed25519"
	"golang.org/x/crypto/blake2b"
)

// Wallet represents a wallet.
type Wallet struct {
	seed  []byte
	index uint32
	key   map[string][]byte
	rpc   rpc.Client
}

// NewWallet creates a new wallet.
func NewWallet(seed []byte) (wallet *Wallet, err error) {
	wallet = &Wallet{
		seed: seed,
		key:  make(map[string][]byte),
		rpc:  rpc.Client{URL: "https://mynano.ninja/api/node"},
	}
	return
}

// NewAccount creates a new account.
func (w *Wallet) NewAccount() (address string, err error) {
	key, err := deriveKey(w.seed, w.index)
	if err != nil {
		return
	}
	pubkey, err := derivePubkey(key)
	if err != nil {
		return
	}
	address, err = deriveAddress(pubkey)
	if err != nil {
		return
	}
	w.key[address] = key
	w.index++
	return
}

// ReceivePendings receives all pending amounts.
func (w *Wallet) ReceivePendings() (err error) {
	accounts := make([]string, 0, len(w.key))
	for address := range w.key {
		accounts = append(accounts, address)
	}
	pendings, err := w.rpc.AccountsPending(accounts, -1)
	for account, pendings := range pendings {
		for s, pending := range pendings {
			blockHash, err := hex.DecodeString(s)
			if err != nil {
				return err
			}
			block := &rpc.Block{
				Type:           "state",
				Account:        account,
				Previous:       make(rpc.BlockHash, 32),
				Representative: "nano_1stofnrxuz3cai7ze75o174bpm7scwj9jn3nxsn8ntzg784jf1gzn1jjdkou",
				Balance:        pending.Amount,
				Link:           blockHash,
				LinkAsAccount:  pending.Source,
			}
			w.signBlock(block)
			blockHash, err = w.rpc.Process(block, "receive")
			if err != nil {
				return err
			}
		}
	}
	return
}

func (w *Wallet) signBlock(block *rpc.Block) (err error) {
	hash, err := blake2b.New256(nil)
	if err != nil {
		return
	}
	hash.Write(make([]byte, 31))
	hash.Write([]byte{6})
	pubkey, err := addressToPubkey(block.Account)
	if err != nil {
		return
	}
	hash.Write(pubkey)
	hash.Write(block.Previous)
	pubkey, err = addressToPubkey(block.Representative)
	if err != nil {
		return
	}
	hash.Write(pubkey)
	hash.Write(block.Balance.FillBytes(make([]byte, 16)))
	hash.Write(block.Link)
	block.Signature = ed25519.Sign(w.key[block.Account], hash.Sum(nil))
	return
}
