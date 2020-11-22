package wallet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/hectorchu/gonano/wallet/bip32"
	"github.com/hectorchu/gonano/wallet/ed25519"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/blake2b"
)

func deriveKey(seed []byte, index uint32) (key []byte, err error) {
	if len(seed) != 32 {
		err = errors.New("seed must be 32 bytes")
		return
	}
	hash, err := blake2b.New256(nil)
	if err != nil {
		return
	}
	hash.Write(seed)
	if err = binary.Write(hash, binary.BigEndian, index); err != nil {
		return
	}
	return hash.Sum(nil), nil
}

func deriveKeypair(key []byte) (pubkey, privkey []byte, err error) {
	return ed25519.GenerateKey(bytes.NewReader(key))
}

func newBip39Seed(mnemonic, password string) (seed []byte, err error) {
	return bip39.NewSeedWithErrorChecking(mnemonic, password)
}

func deriveBip39Key(seed []byte, index uint32) (key []byte, err error) {
	key2, err := bip32.NewMasterKey(seed)
	if err != nil {
		return
	}
	for _, i := range []uint32{44, 165, index} {
		if key2, err = key2.NewChildKey(0x80000000 | i); err != nil {
			return
		}
	}
	return key2.Key, nil
}
