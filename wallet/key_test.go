package wallet

import (
	"encoding/hex"
	"testing"

	"github.com/hectorchu/gonano/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeriveKey(t *testing.T) {
	seed, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	key, err := deriveKey(seed, 1)
	require.Nil(t, err)
	assert.Equal(t, "1495f2d49159cc2eaaaa97ebb42346418e1268aff16d7fca90e6bad6d0965520", hex.EncodeToString(key))
}

func TestBip39(t *testing.T) {
	seed, err := newBip39Seed("edge defense waste choose enrich upon flee junk siren film clown finish "+
		"luggage leader kid quick brick print evidence swap drill paddle truly occur", "some password")
	require.Nil(t, err)
	key, err := deriveBip39Key(seed, 0)
	require.Nil(t, err)
	pubkey, _, err := deriveKeypair(key)
	require.Nil(t, err)
	address, err := util.PubkeyToAddress(pubkey)
	require.Nil(t, err)
	assert.Equal(t, "nano_1pu7p5n3ghq1i1p4rhmek41f5add1uh34xpb94nkbxe8g4a6x1p69emk8y1d", address)
}
