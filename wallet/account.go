package wallet

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/hectorchu/gonano/rpc"
	"github.com/hectorchu/gonano/wallet/ed25519"
	"golang.org/x/crypto/blake2b"
)

// Account represents a wallet account.
type Account struct {
	w           *Wallet
	key, pubkey []byte
	address     string
}

// Address returns the address of the account.
func (a *Account) Address() string {
	return a.address
}

// Send sends an amount to an account.
func (a *Account) Send(account string, amount *big.Int) (hash rpc.BlockHash, err error) {
	link, err := addressToPubkey(account)
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
	if err = a.sign(block); err != nil {
		return
	}
	if block.Work, _, _, err = a.w.RPCWork.WorkGenerate(info.Frontier); err != nil {
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
		info.Representative = "nano_1stofnrxuz3cai7ze75o174bpm7scwj9jn3nxsn8ntzg784jf1gzn1jjdkou"
	}
	for hash, pending := range pendings {
		link, err := hex.DecodeString(hash)
		if err != nil {
			return err
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
		if err = a.sign(block); err != nil {
			return err
		}
		if block.Work, _, _, err = a.w.RPCWork.WorkGenerate(workHash); err != nil {
			return err
		}
		if info.Frontier, err = a.w.RPC.Process(block, "receive"); err != nil {
			return err
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
	if err = a.sign(block); err != nil {
		return
	}
	if block.Work, _, _, err = a.w.RPCWork.WorkGenerate(info.Frontier); err != nil {
		return
	}
	return a.w.RPC.Process(block, "change")
}

func (a *Account) sign(block *rpc.Block) (err error) {
	hash, err := blake2b.New256(nil)
	if err != nil {
		return
	}
	hash.Write(make([]byte, 31))
	hash.Write([]byte{6})
	hash.Write(a.pubkey)
	hash.Write(block.Previous)
	pubkey, err := addressToPubkey(block.Representative)
	if err != nil {
		return
	}
	hash.Write(pubkey)
	hash.Write(block.Balance.FillBytes(make([]byte, 16)))
	hash.Write(block.Link)
	block.Signature = ed25519.Sign(a.key, hash.Sum(nil))
	return
}
