package cmd

import (
	"github.com/spf13/cobra"
)

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive all pending amounts for a wallet or account",
	Run: func(cmd *cobra.Command, args []string) {
		if walletIndex < 0 {
			checkWalletAccount()
			for i, wi := range wallets {
				if _, ok := wi.Accounts[walletAccount]; ok {
					walletIndex = i
					break
				}
			}
			if walletIndex < 0 {
				fatal("account not found in any wallet")
			}
		}
		checkWalletIndex()
		wi := wallets[walletIndex]
		wi.init()
		if walletAccount != "" {
			index, ok := wi.Accounts[walletAccount]
			if !ok {
				fatal("account not found in the specified wallet")
			}
			a, err := wi.w.NewAccount(&index)
			fatalIf(err)
			err = a.ReceivePendings()
			fatalIf(err)
		} else {
			err := wi.w.ReceivePendings()
			fatalIf(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
}
