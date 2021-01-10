package cmd

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/hectorchu/gonano/rpc"
	"github.com/hectorchu/gonano/util"
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

			var accounts []string
			for address := range wallets[walletIndex].Accounts {
				accounts = append(accounts, address)
			}

			sort.Strings(accounts)
			rpcClient := rpc.Client{URL: rpcURL}

			var balanceSum, pendingSum big.Int

			for _, address := range accounts {
				balance, pending, err := rpcClient.AccountBalance(context.Background(), address)
				fatalIf(err)

				balanceSum.Add(&balanceSum, &balance.Int)
				pendingSum.Add(&pendingSum, &pending.Int)

				fmt.Print(address)

				printAmounts(&balance.Int, &pending.Int)
			}

			if len(accounts) > 1 {
				fmt.Print(strings.Repeat(" ", 61), "Sum:")
				printAmounts(&balanceSum, &pendingSum)
			}
		}
	},
}

func printAmounts(balance, pending *big.Int) {
	if balance.Cmp(&big.Int{}) > 0 {
		fmt.Printf(" %s", util.NanoAmount{Raw: balance})
	}

	if pending.Cmp(&big.Int{}) > 0 {
		fmt.Printf(" (+ %s pending)", util.NanoAmount{Raw: pending})
	}

	fmt.Println()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
