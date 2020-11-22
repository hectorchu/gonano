package cmd

import (
	"fmt"

	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/cobra"
)

var ledgerCmd = &cobra.Command{
	Use:   "ledger",
	Short: "Add a Ledger wallet",
	Run: func(cmd *cobra.Command, args []string) {
		w, err := wallet.NewLedgerWallet()
		fatalIf(err)
		wi := &walletInfo{w: w, IsLedger: true}
		wallets = append(wallets, wi)
		wi.initAccounts()
		fmt.Println("Added wallet.")
	},
}

func init() {
	addCmd.AddCommand(ledgerCmd)
}
