package wallet

import (
	"context"

	"github.com/hectorchu/gonano/pow"
)

func (w *Wallet) workGenerate(ctx context.Context, data []byte) (work []byte, err error) {
	// nolint:dogsled // TODO - Return as a struct
	_, _, networkMinimum, _, _, _, err := w.RPC.ActiveDifficulty(ctx)
	if err != nil {
		return
	}

	if work, _, _, err = w.RPCWork.WorkGenerate(ctx, data, networkMinimum); err == nil {
		return
	}

	return pow.Generate(data, networkMinimum)
}

func (w *Wallet) workGenerateReceive(ctx context.Context, data []byte) (work []byte, err error) {
	// nolint:dogsled // TODO - Return as a struct
	_, _, _, _, networkReceiveMinimum, _, err := w.RPC.ActiveDifficulty(ctx)
	if err != nil {
		return
	}

	if work, _, _, err = w.RPCWork.WorkGenerate(ctx, data, networkReceiveMinimum); err == nil {
		return
	}

	return pow.Generate(data, networkReceiveMinimum)
}
