package cmd

import (
	"github.com/spf13/cobra"
)

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive all pending amounts for a wallet or account",
	Run: func(cmd *cobra.Command, args []string) {
		if walletAccount == "" {
			checkWalletIndex()

			wi := wallets[walletIndex]

			wi.init()

			for _, index := range wi.Accounts {
				index := index
				_, err := wi.w.NewAccount(&index)
				fatalIf(err)
			}

			err := wi.w.ReceivePendings(cmd.Context())
			fatalIf(err)
		} else {
			err := getAccount().ReceivePendings(cmd.Context())
			fatalIf(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
}
