package wallet

import (
	"bytes"

	"github.com/hectorchu/gonano/ledger"
	"github.com/hectorchu/gonano/rpc"
	"github.com/hectorchu/gonano/wallet/ed25519"
)

type seedImpl struct{}

func (seedImpl) deriveAccount(a *Account) (err error) {
	var key []byte
	if a.w.isBip39 {
		key, err = deriveBip39Key(a.w.seed, a.index)
	} else {
		key, err = deriveKey(a.w.seed, a.index)
	}
	if err != nil {
		return
	}
	a.pubkey, a.key, err = deriveKeypair(key)
	return
}

func (seedImpl) signBlock(a *Account, block *rpc.Block) (err error) {
	hash, err := block.Hash()
	if err != nil {
		return
	}
	block.Signature = ed25519.Sign(a.key, hash)
	return
}

type ledgerImpl struct{}

func (ledgerImpl) deriveAccount(a *Account) (err error) {
	path := []uint32{44, 165, a.index}
	a.pubkey, _, err = ledger.GetAddress(path)
	return
}

func (ledgerImpl) signBlock(a *Account, block *rpc.Block) (err error) {
	path := []uint32{44, 165, a.index}
	var zero [32]byte
	if !bytes.Equal(block.Previous, zero[:]) {
		bi, err := a.w.RPC.BlockInfo(block.Previous)
		if err != nil {
			return err
		}
		if err = ledger.CacheBlock(path, bi.Contents); err != nil {
			return err
		}
	}
	_, block.Signature, err = ledger.SignBlock(path, block)
	return
}
