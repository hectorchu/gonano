package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"syscall"

	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func readPassword(prompt string) (password []byte) {
	fmt.Print(prompt)
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	fatalIf(err)
	return
}

func deriveKey(password, salt []byte) (key, salt2 []byte, err error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err = rand.Read(salt); err != nil {
			return
		}
	}
	key, err = scrypt.Key(password, salt, 1048576, 8, 1, 32)
	salt2 = salt
	return
}

func encrypt(data, key []byte) (enc []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return
	}
	enc = gcm.Seal(nonce, nonce, data, nil)
	return
}

func decrypt(enc, key []byte) (data []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonce := enc[:gcm.NonceSize()]
	enc = enc[gcm.NonceSize():]
	return gcm.Open(nil, nonce, enc, nil)
}
