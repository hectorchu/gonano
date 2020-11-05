package rpc

import "strconv"

func toUint(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
