package wallet

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"errors"

	"github.com/hectorchu/gonano/wallet/ed25519"
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

func derivePubkey(key []byte) (pubkey []byte, err error) {
	pubkey, _, err = ed25519.GenerateKey(bytes.NewReader(key))
	return
}

func deriveAddress(pubkey []byte) (address string, err error) {
	hash, err := blake2b.New(5, nil)
	if err != nil {
		return
	}
	hash.Write(pubkey)
	var checksum []byte
	for _, b := range hash.Sum(nil) {
		checksum = append([]byte{b}, checksum...)
	}
	pubkey = append([]byte{0, 0, 0}, pubkey...)
	b32 := base32.NewEncoding("13456789abcdefghijkmnopqrstuwxyz")
	return "nano_" + b32.EncodeToString(pubkey)[4:] + b32.EncodeToString(checksum), nil
}
