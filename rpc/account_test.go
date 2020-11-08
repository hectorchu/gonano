package rpc_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/hectorchu/gonano/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAccount = "nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny"

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
	balance, pending, err := getClient().AccountBalance(testAccount)
	require.Nil(t, err)
	assertEqualBig(t, "134000000000000000000000000", &balance.Int)
	assertEqualBig(t, "0", &pending.Int)
}

func TestAccountHistory(t *testing.T) {
	history, previous, err := getClient().AccountHistory(testAccount, 1, nil)
	require.Nil(t, err)
	require.Len(t, history, 1)
	h := history[0]
	assert.Equal(t, "receive", h.Type)
	assert.Equal(t, testAccount, h.Account)
	assertEqualBig(t, "100000000000000000000000000", &h.Amount.Int)
	assert.Equal(t, uint64(1604610080), h.LocalTimestamp)
	assert.Equal(t, uint64(3), h.Height)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", h.Hash)
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", previous)
}

func TestAccountHistoryRaw(t *testing.T) {
	history, previous, err := getClient().AccountHistoryRaw(testAccount, 1, nil)
	require.Nil(t, err)
	require.Len(t, history, 1)
	h := history[0]
	assert.Equal(t, "state", h.Type)
	assert.Equal(t, "nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd", h.Representative)
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", h.Link)
	assertEqualBig(t, "134000000000000000000000000", &h.Balance.Int)
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", h.Previous)
	assert.Equal(t, "receive", h.Subtype)
	assert.Equal(t, testAccount, h.Account)
	assertEqualBig(t, "100000000000000000000000000", &h.Amount.Int)
	assert.Equal(t, uint64(1604610080), h.LocalTimestamp)
	assert.Equal(t, uint64(3), h.Height)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", h.Hash)
	assertEqualBytes(t, "788f7ec074f1854b", h.Work)
	assertEqualBytes(t, "E0F2C0187F87917C28BB989DA516114F64FEEAD307011F73F1A0982B3603A51740279ED5DA4D428C3F0E652A638BB75F790B695F9D23125B54DB3312A7F28100", h.Signature)
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", previous)
}

func TestAccountInfo(t *testing.T) {
	i, err := getClient().AccountInfo(testAccount)
	require.Nil(t, err)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", i.Frontier)
	assertEqualBytes(t, "E6F513D4821F60151DD3C08C078AF3403F59AE44CC7983083E2391A3E1972A8F", i.OpenBlock)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", i.RepresentativeBlock)
	assertEqualBig(t, "134000000000000000000000000", &i.Balance.Int)
	assert.Equal(t, uint64(1604610080), i.ModifiedTimestamp)
	assert.Equal(t, uint64(3), i.BlockCount)
	assert.Equal(t, uint64(2), i.AccountVersion)
	assert.Equal(t, uint64(3), i.ConfirmationHeight)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", i.ConfirmationHeightFrontier)
}

func TestAccountKey(t *testing.T) {
	key, err := getClient().AccountKey(testAccount)
	require.Nil(t, err)
	assertEqualBytes(t, "7D4D6D8A612F199FC559676DEB63BD02891F9E4B3841411E44CAB909FFF01D8B", key)
}

func TestAccountRepresentative(t *testing.T) {
	representative, err := getClient().AccountRepresentative(testAccount)
	require.Nil(t, err)
	assert.Equal(t, "nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd", representative)
}

func TestAccountWeight(t *testing.T) {
	weight, err := getClient().AccountWeight(testAccount)
	require.Nil(t, err)
	assertEqualBig(t, "0", &weight.Int)
}

func TestAccountsBalances(t *testing.T) {
	balances, err := getClient().AccountsBalances([]string{testAccount})
	require.Nil(t, err)
	require.Len(t, balances, 1)
	assertEqualBig(t, "134000000000000000000000000", &balances[testAccount].Balance.Int)
	assertEqualBig(t, "0", &balances[testAccount].Pending.Int)
}

func TestAccountsFrontiers(t *testing.T) {
	frontiers, err := getClient().AccountsFrontiers([]string{testAccount})
	require.Nil(t, err)
	require.Len(t, frontiers, 1)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", frontiers[testAccount])
}

func TestAccountsPending(t *testing.T) {
	pendings, err := getClient().AccountsPending([]string{
		testAccount, "nano_159m8t4iedstzcaacikb9hdkhbcxcqzfbw56dutay8ceqagq9wxpsk9ftfq9"}, 1)
	require.Nil(t, err)
	require.Len(t, pendings, 2)
	assert.Empty(t, pendings[testAccount])
	blocks := pendings["nano_159m8t4iedstzcaacikb9hdkhbcxcqzfbw56dutay8ceqagq9wxpsk9ftfq9"]
	require.Len(t, blocks, 1)
	pending := blocks["96D8422D1CB676EF1B62A313865626A7725C3B9BB5B875601A1460ACF30B5322"]
	assertEqualBig(t, "123000000000000000000000000", &pending.Amount.Int)
	assert.Equal(t, "nano_3kwppxjcggzs65fjh771ch6dbuic3xthsn5wsg6i5537jacw7m493ra8574x", pending.Source)
}

func TestFrontierCount(t *testing.T) {
	count, err := getClient().FrontierCount()
	require.Nil(t, err)
	assert.Greater(t, count, uint64(2000000))
}

func TestRepresentatives(t *testing.T) {
	representatives, err := getClient().Representatives(2)
	require.Nil(t, err)
	require.Len(t, representatives, 3)
	for _, weight := range representatives {
		assert.GreaterOrEqual(t, weight.Cmp(&big.Int{}), 0)
	}
}

func TestRepresentativesOnline(t *testing.T) {
	representatives, err := getClient().RepresentativesOnline()
	require.Nil(t, err)
	require.NotEmpty(t, representatives)
	for _, r := range representatives {
		assert.Greater(t, r.Weight.Cmp(&big.Int{}), 0)
	}
}
