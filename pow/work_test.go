package pow_test

import (
	"encoding/binary"
	"math/rand"
	"strconv"
	"testing"

	"github.com/hectorchu/gonano/pow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/blake2b"
)

func TestGenerateGPU(t *testing.T) {
	data := make([]byte, 32)

	rand.Read(data)

	target, _ := strconv.ParseUint("fffffff800000000", 16, 0)
	work, err := pow.GenerateGPU(data, target)

	require.Nil(t, err)

	hash, _ := blake2b.New(8, nil)
	hash.Write(work)
	hash.Write(data)

	assert.True(t, binary.LittleEndian.Uint64(hash.Sum(nil)) >= target)
}

func TestGenerateCPU(t *testing.T) {
	data := make([]byte, 32)

	rand.Read(data)

	target, _ := strconv.ParseUint("fffffe0000000000", 16, 0)
	work, err := pow.GenerateCPU(data, target)

	require.Nil(t, err)

	hash, _ := blake2b.New(8, nil)
	hash.Write(work)
	hash.Write(data)

	assert.True(t, binary.LittleEndian.Uint64(hash.Sum(nil)) >= target)
}
