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

func deriveKeypair(key []byte) (pubkey, privkey []byte, err error) {
	return ed25519.GenerateKey(bytes.NewReader(key))
}

func deriveAddress(pubkey []byte) (address string, err error) {
	checksum, err := checksum(pubkey)
	if err != nil {
		return
	}
	pubkey = append([]byte{0, 0, 0}, pubkey...)
	b32 := base32.NewEncoding("13456789abcdefghijkmnopqrstuwxyz")
	return "nano_" + b32.EncodeToString(pubkey)[4:] + b32.EncodeToString(checksum), nil
}

func checksum(pubkey []byte) (checksum []byte, err error) {
	hash, err := blake2b.New(5, nil)
	if err != nil {
		return
	}
	hash.Write(pubkey)
	for _, b := range hash.Sum(nil) {
		checksum = append([]byte{b}, checksum...)
	}
	return
}

func addressToPubkey(address string) (pubkey []byte, err error) {
	err = errors.New("invalid address")
	switch len(address) {
	case 64:
		if address[:4] != "xrb_" {
			return
		}
		address = address[4:]
	case 65:
		if address[:5] != "nano_" {
			return
		}
		address = address[5:]
	default:
		return
	}
	b32 := base32.NewEncoding("13456789abcdefghijkmnopqrstuwxyz")
	if pubkey, err = b32.DecodeString("1111" + address[:52]); err != nil {
		return
	}
	pubkey = pubkey[3:]
	checksum, err := checksum(pubkey)
	if err != nil {
		return
	}
	if b32.EncodeToString(checksum) != address[52:] {
		err = errors.New("checksum mismatch")
	}
	return
}
