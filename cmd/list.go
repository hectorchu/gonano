package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List wallets or accounts within a wallet",
}

func init() {
	rootCmd.AddCommand(listCmd)
}
