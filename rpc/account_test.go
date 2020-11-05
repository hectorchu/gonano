package rpc_test

import (
	"math/big"
	"testing"

	"github.com/hectorchu/gonano/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getClient() *rpc.Client {
	return &rpc.Client{URL: "https://mynano.ninja/api/node"}
}

func TestAccountBalance(t *testing.T) {
	balance, pending, err := getClient().AccountBalance("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny")
	require.Nil(t, err)
	var x big.Int
	x.SetString("134000000000000000000000000", 10)
	assert.Equal(t, &x, balance)
	assert.Equal(t, &big.Int{}, pending)
}

func TestAccountHistory(t *testing.T) {
	history, previous, err := getClient().AccountHistory("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny", 1)
	require.Nil(t, err)
	require.Len(t, history, 1)
	h := history[0]
	assert.Equal(t, "receive", h.Type)
	assert.Equal(t, "nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny", h.Account)
	var x big.Int
	x.SetString("100000000000000000000000000", 10)
	assert.Equal(t, &x, h.Amount)
	assert.Equal(t, uint64(1604610080), h.LocalTimestamp)
	assert.Equal(t, uint64(3), h.Height)
	assert.Equal(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", h.Hash)
	assert.Equal(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", previous)
}

func TestAccountKey(t *testing.T) {
	key, err := getClient().AccountKey("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny")
	require.Nil(t, err)
	assert.Equal(t, "7D4D6D8A612F199FC559676DEB63BD02891F9E4B3841411E44CAB909FFF01D8B", key)
}

func TestAccountRepresentative(t *testing.T) {
	representative, err := getClient().AccountRepresentative("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny")
	require.Nil(t, err)
	assert.Equal(t, "nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd", representative)
}

func TestAccountWeight(t *testing.T) {
	weight, err := getClient().AccountWeight("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny")
	require.Nil(t, err)
	assert.Equal(t, &big.Int{}, weight)
}
