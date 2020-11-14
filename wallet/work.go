package wallet

import (
	"encoding/binary"
	"hash"
	"math/rand"
	"runtime"

	"golang.org/x/crypto/blake2b"
)

func (w *Wallet) workGenerate(data []byte) (work []byte, err error) {
	_, _, networkMinimum, _, _, _, err := w.RPC.ActiveDifficulty()
	if err != nil {
		return
	}
	if work, _, _, err = w.RPCWork.WorkGenerate(data, networkMinimum); err == nil {
		return
	}
	return workGenerate(data, networkMinimum)
}

func (w *Wallet) workGenerateReceive(data []byte) (work []byte, err error) {
	_, _, _, _, networkReceiveMinimum, _, err := w.RPC.ActiveDifficulty()
	if err != nil {
		return
	}
	if work, _, _, err = w.RPCWork.WorkGenerate(data, networkReceiveMinimum); err == nil {
		return
	}
	return workGenerate(data, networkReceiveMinimum)
}

func workGenerate(data, difficulty []byte) (work []byte, err error) {
	target := binary.BigEndian.Uint64(difficulty)
	n := runtime.NumCPU()
	ch := make(chan []byte, n)
	hash := make([]hash.Hash, n)
	for i := 0; i < n; i++ {
		if hash[i], err = blake2b.New(8, nil); err != nil {
			return
		}
	}
	done := false
	x := rand.Uint64()
	for i := 0; i < n; i++ {
		go func(i int) {
			work := make([]byte, 8)
			for x := x + uint64(i); !done; x += uint64(n) {
				binary.BigEndian.PutUint64(work, x)
				hash[i].Reset()
				hash[i].Write(work)
				hash[i].Write(data)
				if binary.LittleEndian.Uint64(hash[i].Sum(nil)) >= target {
					done = true
					ch <- work
				}
			}
		}(i)
	}
	work = <-ch
	for i, j := 0, len(work)-1; i < j; i, j = i+1, j-1 {
		work[i], work[j] = work[j], work[i]
	}
	return
}
