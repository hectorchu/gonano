package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List wallets",
	Long:  `List wallets.`,
	Run: func(cmd *cobra.Command, args []string) {
		for i, wi := range wallets {
			fmt.Printf("%d: %d account(s)\n", i, len(wi.Accounts))
		}
	},
}

func init() {
	walletCmd.AddCommand(listCmd)
}
