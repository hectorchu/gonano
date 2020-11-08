package rpc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActiveDifficulty(t *testing.T) {
	multiplier, networkCurrent, networkMinimum,
		networkReceiveCurrent, networkReceiveMinimum,
		difficultyTrend, err := getClient().ActiveDifficulty()
	require.Nil(t, err)
	assert.GreaterOrEqual(t, multiplier, 1.0)
	assert.Len(t, networkCurrent, 8)
	assert.Len(t, networkMinimum, 8)
	assert.Len(t, networkReceiveCurrent, 8)
	assert.Len(t, networkReceiveMinimum, 8)
	assert.Len(t, difficultyTrend, 20)
	for _, multiplier := range difficultyTrend {
		assert.GreaterOrEqual(t, multiplier, 1.0)
	}
}

func TestAvailableSupply(t *testing.T) {
	available, err := getClient().AvailableSupply()
	require.Nil(t, err)
	assertEqualBig(t, "133246497546603000000000000000000000000", &available.Int)
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

func TestFrontierCount(t *testing.T) {
	count, err := getClient().FrontierCount()
	require.Nil(t, err)
	assert.Greater(t, count, uint64(2000000))
}
