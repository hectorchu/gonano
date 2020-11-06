package rpc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"
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

func toTime(x interface{}) (t time.Time, err error) {
	v, err := toUint(x)
	if err != nil {
		return
	}
	t = time.Unix(int64(v), 0).UTC()
	return
}

func toBytes(x interface{}) (b []byte, err error) {
	s, err := toStr(x)
	if err != nil {
		return
	}
	return hex.DecodeString(s)
}
