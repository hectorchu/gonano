package wallet

import (
	"github.com/hectorchu/gonano/rpc"
)

// Wallet represents a wallet.
type Wallet struct {
	seed         []byte
	index        uint32
	accounts     map[string]*Account
	RPC, RPCWork rpc.Client
}

// NewWallet creates a new wallet.
func NewWallet(seed []byte) *Wallet {
	return &Wallet{
		seed:     seed,
		accounts: make(map[string]*Account),
		RPC:      rpc.Client{URL: "https://mynano.ninja/api/node"},
		RPCWork:  rpc.Client{URL: "http://[::1]:7076"},
	}
}

// NewAccount creates a new account.
func (w *Wallet) NewAccount() (a *Account, err error) {
	key, err := deriveKey(w.seed, w.index)
	if err != nil {
		return
	}
	a = &Account{w: w}
	if a.pubkey, a.key, err = deriveKeypair(key); err != nil {
		return
	}
	if a.address, err = deriveAddress(a.pubkey); err != nil {
		return
	}
	w.accounts[a.address] = a
	w.index++
	return
}

// ReceivePendings pockets all pending amounts.
func (w *Wallet) ReceivePendings() (err error) {
	accounts := make([]string, 0, len(w.accounts))
	for address := range w.accounts {
		accounts = append(accounts, address)
	}
	pendings, err := w.RPC.AccountsPending(accounts, -1)
	if err != nil {
		return
	}
	for account, pendings := range pendings {
		if err = w.accounts[account].receivePendings(pendings); err != nil {
			return
		}
	}
	return
}
