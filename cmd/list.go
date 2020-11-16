package cmd

import (
	"fmt"
	"math/big"

	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List wallets or accounts within a wallet",
	Run: func(cmd *cobra.Command, args []string) {
		if walletIndex < 0 {
			for i, wi := range wallets {
				n := len(wi.Accounts)
				switch n {
				case 1:
					fmt.Printf("%d: %d account\n", i, n)
				default:
					fmt.Printf("%d: %d accounts\n", i, n)
				}
			}
		} else {
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
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
