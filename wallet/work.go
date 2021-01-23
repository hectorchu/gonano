package wallet

import "github.com/hectorchu/gonano/pow"

func (w *Wallet) workGenerate(data []byte) (work []byte, err error) {
	_, networkCurrent, _, _, _, _, err := w.RPC.ActiveDifficulty()
	if err != nil {
		return
	}
	if work, _, _, err = w.RPCWork.WorkGenerate(data, networkCurrent); err == nil {
		return
	}
	return pow.Generate(data, networkCurrent)
}

func (w *Wallet) workGenerateReceive(data []byte) (work []byte, err error) {
	_, _, _, networkReceiveCurrent, _, _, err := w.RPC.ActiveDifficulty()
	if err != nil {
		return
	}
	if work, _, _, err = w.RPCWork.WorkGenerate(data, networkReceiveCurrent); err == nil {
		return
	}
	return pow.Generate(data, networkReceiveCurrent)
}
