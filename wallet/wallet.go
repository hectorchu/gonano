package wallet

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/hectorchu/gonano/rpc"
	"github.com/hectorchu/gonano/wallet/ed25519"
	"golang.org/x/crypto/blake2b"
)

// Wallet represents a wallet.
type Wallet struct {
	seed         []byte
	index        uint32
	key          map[string][]byte
	RPC, RPCWork rpc.Client
}

// NewWallet creates a new wallet.
func NewWallet(seed []byte) (wallet *Wallet, err error) {
	wallet = &Wallet{
		seed:    seed,
		key:     make(map[string][]byte),
		RPC:     rpc.Client{URL: "https://mynano.ninja/api/node"},
		RPCWork: rpc.Client{URL: "http://[::1]:7076"},
	}
	return
}

// NewAccount creates a new account.
func (w *Wallet) NewAccount() (address string, err error) {
	key, err := deriveKey(w.seed, w.index)
	if err != nil {
		return
	}
	pubkey, key, err := deriveKeypair(key)
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

// Send sends an amount from an account.
func (w *Wallet) Send(account, destAccount string, amount *big.Int) (hash rpc.BlockHash, err error) {
	info, err := w.RPC.AccountInfo(account)
	if err != nil {
		return
	}
	link, err := addressToPubkey(destAccount)
	if err != nil {
		return
	}
	info.Balance.Sub(&info.Balance.Int, amount)
	if info.Balance.Cmp(&big.Int{}) < 0 {
		err = errors.New("insufficient funds")
		return
	}
	block := &rpc.Block{
		Type:           "state",
		Account:        account,
		Previous:       info.Frontier,
		Representative: info.Representative,
		Balance:        info.Balance,
		Link:           link,
		LinkAsAccount:  destAccount,
	}
	if err = w.sign(block); err != nil {
		return
	}
	if block.Work, _, _, err = w.RPCWork.WorkGenerate(info.Frontier); err != nil {
		return
	}
	return w.RPC.Process(block, "send")
}

// ReceivePendings receives all pending amounts.
func (w *Wallet) ReceivePendings() (err error) {
	accounts := make([]string, 0, len(w.key))
	for address := range w.key {
		accounts = append(accounts, address)
	}
	pendings, err := w.RPC.AccountsPending(accounts, -1)
	if err != nil {
		return
	}
	for account, pendings := range pendings {
		if len(pendings) == 0 {
			continue
		}
		info, err := w.RPC.AccountInfo(account)
		if err != nil {
			info.Balance = &rpc.RawAmount{}
			info.Representative = "nano_1stofnrxuz3cai7ze75o174bpm7scwj9jn3nxsn8ntzg784jf1gzn1jjdkou"
		}
		for hash, pending := range pendings {
			link, err := hex.DecodeString(hash)
			if err != nil {
				return err
			}
			info.Balance.Add(&info.Balance.Int, &pending.Amount.Int)
			if info.Frontier, err = w.receive(
				account, info.Representative, pending.Source,
				info.Balance, info.Frontier, link,
			); err != nil {
				return err
			}
		}
	}
	return
}

func (w *Wallet) receive(
	account, representative, sourceAccount string,
	balance *rpc.RawAmount, previous, link rpc.BlockHash,
) (hash rpc.BlockHash, err error) {
	workHash := previous
	if previous == nil {
		previous = make(rpc.BlockHash, 32)
		if workHash, err = addressToPubkey(account); err != nil {
			return
		}
	}
	block := &rpc.Block{
		Type:           "state",
		Account:        account,
		Previous:       previous,
		Representative: representative,
		Balance:        balance,
		Link:           link,
		LinkAsAccount:  sourceAccount,
	}
	if err = w.sign(block); err != nil {
		return
	}
	if block.Work, _, _, err = w.RPCWork.WorkGenerate(workHash); err != nil {
		return
	}
	return w.RPC.Process(block, "receive")
}

func (w *Wallet) sign(block *rpc.Block) (err error) {
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
