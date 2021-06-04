package rpc_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvailableSupply(t *testing.T) {
	available, err := getClient().AvailableSupply()
	require.Nil(t, err)
	expectedSupply, _ := new(big.Int).SetString("133000000000000000000000000000000000000", 10)
	assert.True(t, available.Cmp(expectedSupply) > 0)
}
