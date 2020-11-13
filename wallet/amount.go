package wallet

import (
	"errors"
	"math/big"
)

// NanoAmount wraps a raw amount.
type NanoAmount struct {
	Raw *big.Int
}

func (NanoAmount) exp() *big.Int {
	x := big.NewInt(10)
	return x.Exp(x, big.NewInt(30), nil)
}

// NanoAmountFromString parses NANO amounts in strings.
func NanoAmountFromString(s string) (n NanoAmount, err error) {
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		err = errors.New("unable to parse nano amount")
		return
	}
	r = r.Mul(r, new(big.Rat).SetInt(n.exp()))
	if !r.IsInt() {
		err = errors.New("unable to parse nano amount")
		return
	}
	n.Raw = r.Num()
	return
}

func (n NanoAmount) String() string {
	r := new(big.Rat).SetFrac(n.Raw, n.exp())
	s := r.FloatString(30)
	return s[:len(s)-24]
}
