package wallet

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/hectorchu/gonano/rpc"
	"github.com/hectorchu/gonano/util"
)

// Account represents a wallet account.
type Account struct {
	w           *Wallet
	index       uint32
	key, pubkey []byte
	address     string
}

// Address returns the address of the account.
func (a *Account) Address() string {
	return a.address
}

// Index returns the derivation index of the account.
func (a *Account) Index() uint32 {
	return a.index
}

// Balance gets the confirmed and pending balances for account.
func (a *Account) Balance() (balance, pending *big.Int, err error) {
	b, p, err := a.w.RPC.AccountBalance(a.address)
	if err != nil {
		return
	}
	return &b.Int, &p.Int, nil
}

// Send sends an amount to an account.
func (a *Account) Send(account string, amount *big.Int) (hash rpc.BlockHash, err error) {
	link, err := util.AddressToPubkey(account)
	if err != nil {
		return
	}
	info, err := a.w.RPC.AccountInfo(a.address)
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
		Account:        a.address,
		Previous:       info.Frontier,
		Representative: info.Representative,
		Balance:        info.Balance,
		Link:           link,
	}
	if err = a.w.impl.signBlock(a, block); err != nil {
		return
	}
	if block.Work, err = a.w.workGenerate(info.Frontier); err != nil {
		return
	}
	return a.w.RPC.Process(block, "send")
}

// ReceivePendings pockets all pending amounts.
func (a *Account) ReceivePendings() (err error) {
	pendings, err := a.w.RPC.AccountsPending([]string{a.address}, -1)
	if err != nil {
		return
	}
	return a.receivePendings(pendings[a.address])
}

func (a *Account) receivePendings(pendings rpc.HashToPendingMap) (err error) {
	if len(pendings) == 0 {
		return
	}
	info, err := a.w.RPC.AccountInfo(a.address)
	if err != nil {
		info.Balance = &rpc.RawAmount{}
		info.Representative = "nano_3gonano8jnse4zm65jaiki9tk8ry4jtgc1smarinukho6fmbc45k3icsh6en"
		err = nil
	}
	for hash, pending := range pendings {
		var link rpc.BlockHash
		if link, err = hex.DecodeString(hash); err != nil {
			return
		}
		workHash := info.Frontier
		if info.Frontier == nil {
			info.Frontier = make(rpc.BlockHash, 32)
			workHash = a.pubkey
		}
		info.Balance.Add(&info.Balance.Int, &pending.Amount.Int)
		block := &rpc.Block{
			Type:           "state",
			Account:        a.address,
			Previous:       info.Frontier,
			Representative: info.Representative,
			Balance:        info.Balance,
			Link:           link,
		}
		if err = a.w.impl.signBlock(a, block); err != nil {
			return
		}
		if block.Work, err = a.w.workGenerateReceive(workHash); err != nil {
			return
		}
		if info.Frontier, err = a.w.RPC.Process(block, "receive"); err != nil {
			return
		}
	}
	return
}

// ChangeRep changes the account's representative.
func (a *Account) ChangeRep(representative string) (hash rpc.BlockHash, err error) {
	info, err := a.w.RPC.AccountInfo(a.address)
	if err != nil {
		return
	}
	block := &rpc.Block{
		Type:           "state",
		Account:        a.address,
		Previous:       info.Frontier,
		Representative: representative,
		Balance:        info.Balance,
		Link:           make(rpc.BlockHash, 32),
	}
	if err = a.w.impl.signBlock(a, block); err != nil {
		return
	}
	if block.Work, err = a.w.workGenerate(info.Frontier); err != nil {
		return
	}
	return a.w.RPC.Process(block, "change")
}
