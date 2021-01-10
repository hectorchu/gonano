package rpc_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActiveDifficulty(t *testing.T) {
	multiplier, networkCurrent, networkMinimum,
		networkReceiveCurrent, networkReceiveMinimum,
		difficultyTrend, err := getClient().ActiveDifficulty(context.Background())

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
	available, err := getClient().AvailableSupply(context.Background())
	require.Nil(t, err)

	expectedSupply, _ := new(big.Int).SetString("133246497000000000000000000000000000000", 10)

	assert.True(t, available.Cmp(expectedSupply) > 0)
}
