package util

import (
	"encoding/base32"
	"errors"

	"golang.org/x/crypto/blake2b"
)

// AddressToPubkey converts address to a pubkey.
func AddressToPubkey(address string) (pubkey []byte, err error) {
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

// PubkeyToAddress converts pubkey to an address.
func PubkeyToAddress(pubkey []byte) (address string, err error) {
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
