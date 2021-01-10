// Package ledger handles ledger hardware wallet functionality
package ledger

import (
	"encoding/binary"

	"github.com/hectorchu/gonano/rpc"
	"github.com/hectorchu/gonano/util"
)

// nolint:gomnd // TODO - Resolve magic numbers
func appendPath(buf []byte, path []uint32) []byte {
	var p [4]byte

	buf = append(buf, byte(len(path)))

	for _, i := range path {
		binary.BigEndian.PutUint32(p[:], 0x80000000|i)
		buf = append(buf, p[:]...)
	}

	return buf
}

// GetAddress returns the public key and the encoded address for the given BIP32 path.
// nolint:gomnd // TODO - Resolve magic numbers
func GetAddress(path []uint32) (pubkey []byte, address string, err error) {
	d, err := getDevice()
	if err != nil {
		return
	}

	defer func() { _ = d.Close() }()

	payload := []byte{0xa1, 0x02, 0x00, 0x00, 0x00}
	payload = appendPath(payload, path)
	payload[4] = byte(len(payload)) - 5

	resp, err := send(d, payload)
	if err != nil {
		return
	}

	pubkey = resp[:32]
	address = string(resp[33 : 33+resp[32]])

	return
}

// CacheBlock caches the frontier block in memory. The sign block payload uses this
// cached data to determine the changes in account state.
// nolint:gomnd // TODO - Resolve magic numbers
func CacheBlock(path []uint32, block *rpc.Block) error {
	d, err := getDevice()
	if err != nil {
		return err
	}

	defer func() { _ = d.Close() }()

	payload := []byte{0xa1, 0x03, 0x00, 0x00, 0x00}
	payload = appendPath(payload, path)
	payload = append(payload, block.Previous...)
	payload = append(payload, block.Link...)

	pubkey, err := util.AddressToPubkey(block.Representative)
	if err != nil {
		return err
	}

	payload = append(payload, pubkey...)
	payload = append(payload, block.Balance.FillBytes(make([]byte, 16))...)
	payload = append(payload, block.Signature...)
	payload[4] = byte(len(payload)) - 5

	_, err = send(d, payload)

	return err
}

// SignBlock returns the signature for the provided block data. For non-null
// parent blocks the cache block payload needs to be called before this payload.
// nolint:gomnd // TODO - Resolve magic numbers
func SignBlock(path []uint32, block *rpc.Block) (rpc.BlockHash, []byte, error) {
	d, err := getDevice()
	if err != nil {
		return nil, nil, err
	}

	defer func() { _ = d.Close() }()

	payload := []byte{0xa1, 0x04, 0x00, 0x00, 0x00}
	payload = appendPath(payload, path)
	payload = append(payload, block.Previous...)
	payload = append(payload, block.Link...)

	pubkey, err := util.AddressToPubkey(block.Representative)
	if err != nil {
		return nil, nil, err
	}

	payload = append(payload, pubkey...)
	payload = append(payload, block.Balance.FillBytes(make([]byte, 16))...)
	payload[4] = byte(len(payload)) - 5

	resp, err := send(d, payload)
	if err != nil {
		return nil, nil, err
	}

	hash := resp[:32]
	signature := resp[32:96]

	return hash, signature, err
}
