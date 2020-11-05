package rpc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockCount(t *testing.T) {
	cemented, count, unchecked, err := getClient().BlockCount()
	require.Nil(t, err)
	assert.Greater(t, cemented, uint64(50000000))
	assert.Greater(t, count, uint64(50000000))
	assert.Less(t, unchecked, uint64(100000))
}
