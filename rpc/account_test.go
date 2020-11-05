package rpc_test

import (
	"math/big"
	"testing"
	"time"

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
	history, previous, err := getClient().AccountHistory("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny", 1, "")
	require.Nil(t, err)
	require.Len(t, history, 1)
	h := history[0]
	assert.Equal(t, "receive", h.Type)
	assert.Equal(t, "nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny", h.Account)
	var x big.Int
	x.SetString("100000000000000000000000000", 10)
	assert.Equal(t, &x, h.Amount)
	assert.Equal(t, time.Date(2020, time.November, 5, 21, 1, 20, 0, time.UTC), h.LocalTimestamp)
	assert.Equal(t, uint64(3), h.Height)
	assert.Equal(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", h.Hash)
	assert.Equal(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", previous)
}

func TestAccountInfo(t *testing.T) {
	i, err := getClient().AccountInfo("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny")
	require.Nil(t, err)
	assert.Equal(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", i.Frontier)
	assert.Equal(t, "E6F513D4821F60151DD3C08C078AF3403F59AE44CC7983083E2391A3E1972A8F", i.OpenBlock)
	assert.Equal(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", i.RepresentativeBlock)
	var x big.Int
	x.SetString("134000000000000000000000000", 10)
	assert.Equal(t, &x, i.Balance)
	assert.Equal(t, time.Date(2020, time.November, 5, 21, 1, 20, 0, time.UTC), i.ModifiedTimestamp)
	assert.Equal(t, uint64(3), i.BlockCount)
	assert.Equal(t, uint64(2), i.AccountVersion)
	assert.Equal(t, uint64(3), i.ConfirmationHeight)
	assert.Equal(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", i.ConfirmationHeightFrontier)
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
