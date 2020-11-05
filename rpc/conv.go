package rpc

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

func toStr(x interface{}) (s string, err error) {
	s, ok := x.(string)
	if !ok {
		err = fmt.Errorf("failed to cast to string: %v", x)
	}
	return
}

func toUint(x interface{}) (v uint64, err error) {
	s, err := toStr(x)
	if err != nil {
		return
	}
	v, err = strconv.ParseUint(s, 10, 64)
	return
}

func toBig(x interface{}) (z *big.Int, err error) {
	s, err := toStr(x)
	if err != nil {
		return
	}
	z = new(big.Int)
	if _, ok := z.SetString(s, 10); !ok {
		err = errors.New("failed to parse big number: " + s)
	}
	return
}
