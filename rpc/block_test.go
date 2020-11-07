package rpc_test

import (
	"encoding/hex"
	"testing"

	"github.com/hectorchu/gonano/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func hexString(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func strToRaw(s string) *rpc.RawAmount {
	var r rpc.RawAmount
	r.SetString(s, 10)
	return &r
}

func TestBlockAccount(t *testing.T) {
	hash := hexString("8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD")
	account, err := getClient().BlockAccount(hash)
	require.Nil(t, err)
	assert.Equal(t, testAccount, account)
}

func TestBlockCount(t *testing.T) {
	cemented, count, unchecked, err := getClient().BlockCount()
	require.Nil(t, err)
	assert.Greater(t, cemented, uint64(50000000))
	assert.Greater(t, count, uint64(50000000))
	assert.Less(t, unchecked, uint64(100000))
}

func TestBlockCountType(t *testing.T) {
	send, receive, open, change, state, err := getClient().BlockCountType()
	require.Nil(t, err)
	assert.Greater(t, send, uint64(5000000))
	assert.Greater(t, receive, uint64(4000000))
	assert.Greater(t, open, uint64(500000))
	assert.Greater(t, change, uint64(20000))
	assert.Greater(t, state, uint64(40000000))
}

func TestBlockHash(t *testing.T) {
	hash, err := getClient().BlockHash(&rpc.Block{
		Type:           "state",
		Account:        "nano_3qgmh14nwztqw4wmcdzy4xpqeejey68chx6nciczwn9abji7ihhum9qtpmdr",
		Previous:       hexString("F47B23107E5F34B2CE06F562B5C435DF72A533251CB414C51B2B62A8F63A00E4"),
		Representative: "nano_1hza3f7wiiqa7ig3jczyxj5yo86yegcmqk3criaz838j91sxcckpfhbhhra1",
		Balance:        strToRaw("1000000000000000000000"),
		Link:           hexString("19D3D919475DEED4696B5D13018151D1AF88B2BD3BCFF048B45031C1F36D1858"),
		LinkAsAccount:  "nano_18gmu6engqhgtjnppqam181o5nfhj4sdtgyhy36dan3jr9spt84rzwmktafc",
		Signature:      hexString("3BFBA64A775550E6D49DF1EB8EEC2136DCD74F090E2ED658FBD9E80F17CB1C9F9F7BDE2B93D95558EC2F277FFF15FD11E6E2162A1714731B743D1E941FA4560A"),
		Work:           hexString("cab7404f0b5449d0"),
	})
	require.Nil(t, err)
	assertEqualBytes(t, "FF0144381CFF0B2C079A115E7ADA7E96F43FD219446E7524C48D1CC9900C4F17", hash)
}
