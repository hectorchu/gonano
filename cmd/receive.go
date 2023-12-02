package cmd

import (
	"github.com/spf13/cobra"
)

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive all receivable amounts for a wallet or account",
	Run: func(cmd *cobra.Command, args []string) {
		if walletAccount == "" {
			checkWalletIndex()
			wi := wallets[walletIndex]
			wi.init()
			for _, index := range wi.Accounts {
				_, err := wi.w.NewAccount(&index)
				fatalIf(err)
			}
			err := wi.w.ReceiveReceivables()
			fatalIf(err)
		} else {
			err := getAccount().ReceiveReceivables()
			fatalIf(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
}
