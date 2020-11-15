package cmd

import (
	"github.com/spf13/cobra"
)

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive all pending amounts for a wallet",
	Run: func(cmd *cobra.Command, args []string) {
		checkWalletIndex()
		wi := wallets[walletIndex]
		wi.init()
		err := wi.w.ReceivePendings()
		fatalIf(err)
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
}
