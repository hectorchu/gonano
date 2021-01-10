// Package util provides common nano util functions
package util

import (
	"encoding/base32"
	"errors"

	"github.com/hectorchu/gonano/constants"
	"golang.org/x/crypto/blake2b"
)

// AddressToPubkey converts address to a pubkey.
func AddressToPubkey(address string) ([]byte, error) {
	const (
		xrbLength  = 64
		nanoLength = 65
	)

	err := errors.New("invalid address")

	switch len(address) {
	case xrbLength:
		if address[:4] != "xrb_" {
			return nil, err
		}

		address = address[4:]
	case nanoLength:
		if address[:5] != "nano_" {
			return nil, err
		}

		address = address[5:]
	default:
		return nil, err
	}

	b32 := base32.NewEncoding("13456789abcdefghijkmnopqrstuwxyz")

	pubkey, err := b32.DecodeString("1111" + address[:52])
	if err != nil {
		return nil, err
	}

	pubkey = pubkey[3:]

	checksum, err := checksum(pubkey)
	if err != nil {
		return nil, err
	}

	if b32.EncodeToString(checksum) != address[52:] {
		return nil, errors.New("checksum mismatch")
	}

	return pubkey, nil
}

// PubkeyToAddress converts pubkey to an address.
func PubkeyToAddress(pubkey []byte) (address string, err error) {
	if len(pubkey) != constants.KeyLength {
		return "", errors.New("invalid pubkey length")
	}

	checksum, err := checksum(pubkey)
	if err != nil {
		return
	}

	pubkey = append([]byte{0, 0, 0}, pubkey...)
	b32 := base32.NewEncoding("13456789abcdefghijkmnopqrstuwxyz")

	return "nano_" + b32.EncodeToString(pubkey)[4:] + b32.EncodeToString(checksum), nil
}

func checksum(pubkey []byte) ([]byte, error) {
	hash, err := blake2b.New(constants.CheckSumLength, nil)
	if err != nil {
		return nil, err
	}

	if _, err := hash.Write(pubkey); err != nil {
		return nil, err
	}

	checksum := []byte{}
	for _, b := range hash.Sum(nil) {
		checksum = append([]byte{b}, checksum...)
	}

	return checksum, nil
}
