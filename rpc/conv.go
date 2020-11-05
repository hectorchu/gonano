package rpc

import (
	"errors"
	"math/big"
	"strconv"
)

func toUint(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func toBig(s string) (*big.Int, error) {
	var x big.Int
	if _, ok := x.SetString(s, 10); !ok {
		return nil, errors.New("failed to parse big number: " + s)
	}
	return &x, nil
}
