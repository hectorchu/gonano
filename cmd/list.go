package cmd

import (
	"fmt"
	"math/big"

	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List wallets",
	Long:  `List wallets.`,
	Run: func(cmd *cobra.Command, args []string) {
		if walletIndex < 0 {
			for i, wi := range wallets {
				fmt.Printf("%d: %d account(s)\n", i, len(wi.Accounts))
			}
			return
		}
		checkWalletIndex()
		for _, account := range wallets[walletIndex].Accounts {
			balance, pending, err := rpcClient.AccountBalance(account)
			fatalIf(err)
			fmt.Print(account)
			if balance.Int.Cmp(&big.Int{}) > 0 {
				fmt.Printf(" %s", wallet.NanoAmount{Raw: &balance.Int})
			}
			if pending.Int.Cmp(&big.Int{}) > 0 {
				fmt.Printf(" (+ %s pending)", wallet.NanoAmount{Raw: &pending.Int})
			}
			fmt.Println()
		}
	},
}

func init() {
	walletCmd.AddCommand(listCmd)
}
