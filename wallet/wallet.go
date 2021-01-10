// Package wallet provides nano wallet functionality
package wallet

import (
	"context"
	"math/big"

	"github.com/hectorchu/gonano/rpc"
	"github.com/hectorchu/gonano/util"
)

// Wallet represents a wallet.
type Wallet struct {
	seed         []byte
	isBip39      bool
	nextIndex    uint32
	accounts     map[string]*Account
	RPC, RPCWork rpc.Client
	impl         interface {
		deriveAccount(*Account) error
		signBlock(context.Context, *Account, *rpc.Block) error
	}
}

// NewWallet creates a new wallet.
func NewWallet(seed []byte) (w *Wallet, err error) {
	w = newWallet(seed)
	return
}

// NewBip39Wallet creates a new BIP39 wallet.
func NewBip39Wallet(mnemonic, password string) (w *Wallet, err error) {
	seed, err := newBip39Seed(mnemonic, password)
	if err != nil {
		return
	}

	w = newWallet(seed)
	w.isBip39 = true

	return
}

// NewLedgerWallet creates a new Ledger wallet.
func NewLedgerWallet() (w *Wallet, err error) {
	w = newWallet(nil)
	w.impl = ledgerImpl{}

	return
}

func newWallet(seed []byte) *Wallet {
	return &Wallet{
		seed:     seed,
		accounts: make(map[string]*Account),
		RPC:      rpc.Client{URL: "https://mynano.ninja/api/node"},
		RPCWork:  rpc.Client{URL: "http://[::1]:7076"},
		impl:     seedImpl{},
	}
}

func (w *Wallet) populateAccounts(accounts []string) error {
	for i := range accounts {
		a, err := w.NewAccount(nil)
		if err != nil {
			return err
		}

		accounts[i] = a.Address()
	}

	return nil
}

// ScanForAccounts scans for accounts.
func (w *Wallet) ScanForAccounts(ctx context.Context) error {
	accounts := make([]string, 10)

	err := w.populateAccounts(accounts)
	if err != nil {
		return err
	}

	balances, err := w.RPC.AccountsBalances(ctx, accounts)
	if err != nil {
		return err
	}

	frontiers, err := w.RPC.AccountsFrontiers(ctx, accounts)
	if err != nil {
		return err
	}

	i := len(accounts) - 1
	for ; i >= 0; i-- {
		if balances[accounts[i]].Pending.Cmp(&big.Int{}) > 0 {
			break
		}

		if frontiers[accounts[i]] != nil {
			break
		}

		w.nextIndex = w.accounts[accounts[i]].index
		delete(w.accounts, accounts[i])
	}

	const minAccountLength = 5
	if i < minAccountLength {
		return nil
	}

	return w.ScanForAccounts(ctx)
}

// NewAccount creates a new account.
func (w *Wallet) NewAccount(index *uint32) (a *Account, err error) {
	index2 := w.nextIndex

	if index != nil {
		index2 = *index
	}

	a = &Account{w: w, index: index2}
	if err = w.impl.deriveAccount(a); err != nil {
		return
	}

	if a.address, err = util.PubkeyToAddress(a.pubkey); err != nil {
		return
	}

	if index == nil {
		w.nextIndex++
	}

	if _, ok := w.accounts[a.address]; !ok {
		w.accounts[a.address] = a
	} else if index == nil {
		return w.NewAccount(nil)
	}

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
func (w *Wallet) ReceivePendings(ctx context.Context) (err error) {
	accounts := make([]string, 0, len(w.accounts))
	for address := range w.accounts {
		accounts = append(accounts, address)
	}

	pendings, err := w.RPC.AccountsPending(ctx, accounts, -1)
	if err != nil {
		return
	}

	for account, pendings := range pendings {
		if err = w.accounts[account].receivePendings(ctx, pendings); err != nil {
			return
		}
	}

	return
}
