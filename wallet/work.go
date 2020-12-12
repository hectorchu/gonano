package wallet

import "github.com/hectorchu/gonano/pow"

func (w *Wallet) workGenerate(data []byte) (work []byte, err error) {
	_, _, networkMinimum, _, _, _, err := w.RPC.ActiveDifficulty()
	if err != nil {
		return
	}
	if work, _, _, err = w.RPCWork.WorkGenerate(data, networkMinimum); err == nil {
		return
	}
	return pow.Generate(data, networkMinimum)
}

func (w *Wallet) workGenerateReceive(data []byte) (work []byte, err error) {
	_, _, _, _, networkReceiveMinimum, _, err := w.RPC.ActiveDifficulty()
	if err != nil {
		return
	}
	if work, _, _, err = w.RPCWork.WorkGenerate(data, networkReceiveMinimum); err == nil {
		return
	}
	return pow.Generate(data, networkReceiveMinimum)
}
