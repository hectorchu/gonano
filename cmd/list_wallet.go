package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listWalletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "List wallets",
	Run: func(cmd *cobra.Command, args []string) {
		for i, wi := range wallets {
			n := len(wi.Accounts)
			switch n {
			case 1:
				fmt.Printf("%d: %d account\n", i, n)
			default:
				fmt.Printf("%d: %d accounts\n", i, n)
			}
		}
	},
}

func init() {
	listCmd.AddCommand(listWalletCmd)
}
