package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rescanCmd = &cobra.Command{
	Use:   "rescan",
	Short: "Rescan a wallet for accounts",
	Run: func(cmd *cobra.Command, args []string) {
		checkWalletIndex()
		wi := wallets[walletIndex]
		n := len(wi.Accounts)
		wi.init()
		wi.initAccounts()
		n = len(wi.Accounts) - n
		switch n {
		case 1:
			fmt.Printf("Added %d account.\n", n)
		default:
			fmt.Printf("Added %d accounts.\n", n)
		}
	},
}

func init() {
	rootCmd.AddCommand(rescanCmd)
}
