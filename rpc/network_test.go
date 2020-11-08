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
