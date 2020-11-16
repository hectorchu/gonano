package cmd

import (
	"fmt"
	"math/big"

	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/cobra"
)

var listAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "List accounts within a wallet",
	Run: func(cmd *cobra.Command, args []string) {
		checkWalletIndex()
		for address := range wallets[walletIndex].Accounts {
			balance, pending, err := rpcClient.AccountBalance(address)
			fatalIf(err)
			fmt.Print(address)
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
	listCmd.AddCommand(listAccountCmd)
}
