package rpc_test

import (
	"encoding/hex"
	"strings"
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
	assert.Less(t, unchecked, uint64(500000))
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

const testBlockInfoHash = "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD"

func testBlockInfo(t *testing.T, info *rpc.BlockInfo) {
	assert.Equal(t, "nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny", info.BlockAccount)
	assertEqualBig(t, "100000000000000000000000000", &info.Amount.Int)
	assertEqualBig(t, "134000000000000000000000000", &info.Balance.Int)
	assert.Equal(t, uint64(3), info.Height)
	assert.Equal(t, uint64(1604610080), info.LocalTimestamp)
	assert.Equal(t, true, info.Confirmed)
	assert.Equal(t, "state", info.Contents.Type)
	assert.Equal(t, "nano_1zcffp784drsmz4oksufxfjut1nb5yh6pg43a6h6bkos39zz19ed6a4r36ny", info.Contents.Account)
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", info.Contents.Previous)
	assert.Equal(t, "nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd", info.Contents.Representative)
	assertEqualBig(t, "134000000000000000000000000", &info.Contents.Balance.Int)
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", info.Contents.Link)
	assert.Equal(t, "nano_3mp773x13xf73rati5p5nry7gqbqqzcuop5usefqwjpushjh5u3yat7bzkoj", info.Contents.LinkAsAccount)
	assertEqualBytes(t, "E0F2C0187F87917C28BB989DA516114F64FEEAD307011F73F1A0982B3603A51740279ED5DA4D428C3F0E652A638BB75F790B695F9D23125B54DB3312A7F28100", info.Contents.Signature)
	assertEqualBytes(t, "788f7ec074f1854b", info.Contents.Work)
	assert.Equal(t, "receive", info.Subtype)
}

func TestBlockInfo(t *testing.T) {
	info, err := getClient().BlockInfo(hexString(testBlockInfoHash))
	require.Nil(t, err)
	testBlockInfo(t, &info)
}

func TestBlocksInfo(t *testing.T) {
	blocks, err := getClient().BlocksInfo([]rpc.BlockHash{hexString(testBlockInfoHash)})
	require.Nil(t, err)
	assert.Len(t, blocks, 1)
	testBlockInfo(t, blocks[strings.ToLower(testBlockInfoHash)])
}

func TestChain(t *testing.T) {
	block := hexString("8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD")
	blocks, err := getClient().Chain(block, -1)
	require.Nil(t, err)
	assert.Len(t, blocks, 3)
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", blocks[0])
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", blocks[1])
	assertEqualBytes(t, "E6F513D4821F60151DD3C08C078AF3403F59AE44CC7983083E2391A3E1972A8F", blocks[2])
}

func TestSuccessors(t *testing.T) {
	block := hexString("E6F513D4821F60151DD3C08C078AF3403F59AE44CC7983083E2391A3E1972A8F")
	blocks, err := getClient().Successors(block, -1)
	require.Nil(t, err)
	assert.Len(t, blocks, 3)
	assertEqualBytes(t, "E6F513D4821F60151DD3C08C078AF3403F59AE44CC7983083E2391A3E1972A8F", blocks[0])
	assertEqualBytes(t, "CEC5287A00F5A50E11A80EC3A63C575D37BFD5BAD87BCB1B7E46DBCBE2F1EC3E", blocks[1])
	assertEqualBytes(t, "8C1B5D4BBE27F05C7A888D1E691A07C550A81AFEE16D913EE21E1093888B82FD", blocks[2])
}
