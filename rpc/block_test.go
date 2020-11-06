package rpc_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockAccount(t *testing.T) {
	hash, _ := hex.DecodeString("8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD")
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
