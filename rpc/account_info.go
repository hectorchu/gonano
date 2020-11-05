package rpc

import "math/big"

// AccountInfo returns frontier, open block, change representative block,
// balance, last modified timestamp from local database & block count for
// account.
type AccountInfo struct {
	Frontier                   string
	OpenBlock                  string
	RepresentativeBlock        string
	Balance                    *big.Int
	ModifiedTimestamp          uint64
	BlockCount                 uint64
	ConfirmationHeight         uint64
	ConfirmationHeightFrontier string
	AccountVersion             uint64
}

func (i *AccountInfo) parse(x map[string]interface{}) (err error) {
	if i.Frontier, err = toStr(x["frontier"]); err != nil {
		return
	}
	if i.OpenBlock, err = toStr(x["open_block"]); err != nil {
		return
	}
	if i.RepresentativeBlock, err = toStr(x["representative_block"]); err != nil {
		return
	}
	if i.Balance, err = toBig(x["balance"]); err != nil {
		return
	}
	if i.ModifiedTimestamp, err = toUint(x["modified_timestamp"]); err != nil {
		return
	}
	if i.BlockCount, err = toUint(x["block_count"]); err != nil {
		return
	}
	if i.ConfirmationHeight, err = toUint(x["confirmation_height"]); err != nil {
		return
	}
	if i.ConfirmationHeightFrontier, err = toStr(x["confirmation_height_frontier"]); err != nil {
		return
	}
	if i.AccountVersion, err = toUint(x["account_version"]); err != nil {
		return
	}
	return
}
