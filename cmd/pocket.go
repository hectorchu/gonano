package cmd

import (
	"github.com/spf13/cobra"
)

// pocketCmd represents the pocket command
var pocketCmd = &cobra.Command{
	Use:   "pocket",
	Short: "Pocket pending amounts",
	Long:  `Pocket all pending amounts for a wallet.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkWalletIndex()
		wi := wallets[walletIndex]
		wi.init()
		err := wi.w.ReceivePendings()
		fatalIf(err)
	},
}

func init() {
	walletCmd.AddCommand(pocketCmd)
}
