package rpc_test

import (
	"encoding/hex"
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

func assertEqualBig(t *testing.T, s string, z *big.Int) {
	var x big.Int
	x.SetString(s, 10)
	assert.Equal(t, &x, z)
}

func assertEqualBytes(t *testing.T, s string, b []byte) {
	a, _ := hex.DecodeString(s)
	assert.Equal(t, a, b)
}

func TestAccountBalance(t *testing.T) {
	balance, pending, err := getClient().AccountBalance("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny")
	require.Nil(t, err)
	assertEqualBig(t, "134000000000000000000000000", balance)
	assertEqualBig(t, "0", pending)
}

func TestAccountHistory(t *testing.T) {
	history, previous, err := getClient().AccountHistory("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny", 1, "")
	require.Nil(t, err)
	require.Len(t, history, 1)
	h := history[0]
	assert.Equal(t, "receive", h.Type)
	assert.Equal(t, "nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny", h.Account)
	assertEqualBig(t, "100000000000000000000000000", h.Amount)
	assert.Equal(t, time.Date(2020, time.November, 5, 21, 1, 20, 0, time.UTC), h.LocalTimestamp)
	assert.Equal(t, uint64(3), h.Height)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", h.Hash)
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", previous)
}

func TestAccountHistoryRaw(t *testing.T) {
	history, previous, err := getClient().AccountHistoryRaw("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny", 1, "")
	require.Nil(t, err)
	require.Len(t, history, 1)
	h := history[0]
	assert.Equal(t, "state", h.Type)
	assert.Equal(t, "nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd", h.Representative)
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", h.Link)
	assertEqualBig(t, "134000000000000000000000000", h.Balance)
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", h.Previous)
	assert.Equal(t, "receive", h.Subtype)
	assert.Equal(t, "nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny", h.Account)
	assertEqualBig(t, "100000000000000000000000000", h.Amount)
	assert.Equal(t, time.Date(2020, time.November, 5, 21, 1, 20, 0, time.UTC), h.LocalTimestamp)
	assert.Equal(t, uint64(3), h.Height)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", h.Hash)
	assertEqualBytes(t, "788f7ec074f1854b", h.Work)
	assertEqualBytes(t, "E0F2C0187F87917C28BB989DA516114F64FEEAD307011F73F1A0982B3603A51740279ED5DA4D428C3F0E652A638BB75F790B695F9D23125B54DB3312A7F28100", h.Signature)
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", previous)
}

func TestAccountInfo(t *testing.T) {
	i, err := getClient().AccountInfo("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny")
	require.Nil(t, err)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", i.Frontier)
	assertEqualBytes(t, "E6F513D4821F60151DD3C08C078AF3403F59AE44CC7983083E2391A3E1972A8F", i.OpenBlock)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", i.RepresentativeBlock)
	assertEqualBig(t, "134000000000000000000000000", i.Balance)
	assert.Equal(t, time.Date(2020, time.November, 5, 21, 1, 20, 0, time.UTC), i.ModifiedTimestamp)
	assert.Equal(t, uint64(3), i.BlockCount)
	assert.Equal(t, uint64(2), i.AccountVersion)
	assert.Equal(t, uint64(3), i.ConfirmationHeight)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", i.ConfirmationHeightFrontier)
}

func TestAccountKey(t *testing.T) {
	key, err := getClient().AccountKey("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny")
	require.Nil(t, err)
	assertEqualBytes(t, "7D4D6D8A612F199FC559676DEB63BD02891F9E4B3841411E44CAB909FFF01D8B", key)
}

func TestAccountRepresentative(t *testing.T) {
	representative, err := getClient().AccountRepresentative("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny")
	require.Nil(t, err)
	assert.Equal(t, "nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd", representative)
}

func TestAccountWeight(t *testing.T) {
	weight, err := getClient().AccountWeight("nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny")
	require.Nil(t, err)
	assertEqualBig(t, "0", weight)
}
