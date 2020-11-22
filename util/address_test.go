package util_test

import (
	"encoding/hex"
	"testing"

	"github.com/hectorchu/gonano/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressToPubkey(t *testing.T) {
	pubkey, err := util.AddressToPubkey("nano_1e5aqegc1jb7qe964u4adzmcezyo6o146zb8hm6dft8tkp79za3sxwjym5rx")
	require.Nil(t, err)
	assert.Equal(t, "3068bb1ca04525bb0e416c485fe6a67fd52540227d267cc8b6e8da958a7fa039", hex.EncodeToString(pubkey))
}

func TestPubkeyToAddress(t *testing.T) {
	pubkey, _ := hex.DecodeString("3068bb1ca04525bb0e416c485fe6a67fd52540227d267cc8b6e8da958a7fa039")
	address, err := util.PubkeyToAddress(pubkey)
	require.Nil(t, err)
	assert.Equal(t, "nano_1e5aqegc1jb7qe964u4adzmcezyo6o146zb8hm6dft8tkp79za3sxwjym5rx", address)
}
