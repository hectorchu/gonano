package wallet

import (
	"github.com/hectorchu/gonano/rpc"
)

// Wallet represents a wallet.
type Wallet struct {
	seed         []byte
	accounts     map[string]*Account
	RPC, RPCWork rpc.Client
}

// NewWallet creates a new wallet.
func NewWallet(seed []byte) (w *Wallet, err error) {
	w = &Wallet{
		seed:     seed,
		accounts: make(map[string]*Account),
		RPC:      rpc.Client{URL: "https://mynano.ninja/api/node"},
		RPCWork:  rpc.Client{URL: "http://[::1]:7076"},
	}
	err = w.scanForAccounts()
	return
}

func (w *Wallet) scanForAccounts() (err error) {
	accounts := make([]string, 10)
	for i := range accounts {
		account, err := w.NewAccount()
		if err != nil {
			return err
		}
		accounts[i] = account.Address()
	}
	frontiers, err := w.RPC.AccountsFrontiers(accounts)
	if err != nil {
		return
	}
	i := len(accounts) - 1
	for ; i >= 0; i-- {
		if frontiers[accounts[i]] != nil {
			break
		}
		delete(w.accounts, accounts[i])
	}
	if i < 5 {
		return
	}
	return w.scanForAccounts()
}

// NewAccount creates a new account.
func (w *Wallet) NewAccount() (a *Account, err error) {
	key, err := deriveKey(w.seed, uint32(len(w.accounts)))
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
	return
}

// GetAccount gets the account with address or nil if not found.
func (w *Wallet) GetAccount(address string) *Account {
	return w.accounts[address]
}

// GetAccounts gets all the accounts in the wallet.
func (w *Wallet) GetAccounts() (accounts []*Account) {
	accounts = make([]*Account, 0, len(w.accounts))
	for _, account := range w.accounts {
		accounts = append(accounts, account)
	}
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
