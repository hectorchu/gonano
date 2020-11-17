package cmd

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/hectorchu/gonano/rpc"
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
			rpcClient := rpc.Client{URL: rpcURL}
			var balanceSum, pendingSum big.Int
			for address := range wallets[walletIndex].Accounts {
				balance, pending, err := rpcClient.AccountBalance(address)
				fatalIf(err)
				balanceSum.Add(&balanceSum, &balance.Int)
				pendingSum.Add(&pendingSum, &pending.Int)
				fmt.Print(address)
				printAmounts(&balance.Int, &pending.Int)
			}
			if len(wallets[walletIndex].Accounts) > 1 {
				fmt.Print(strings.Repeat(" ", 61), "Sum:")
				printAmounts(&balanceSum, &pendingSum)
			}
		}
	},
}

func printAmounts(balance, pending *big.Int) {
	if balance.Cmp(&big.Int{}) > 0 {
		fmt.Printf(" %s", wallet.NanoAmount{Raw: balance})
	}
	if pending.Cmp(&big.Int{}) > 0 {
		fmt.Printf(" (+ %s pending)", wallet.NanoAmount{Raw: pending})
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
