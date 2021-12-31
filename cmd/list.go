package cmd

import (
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
		if walletAccount != "" {
			// For when a specific account is specified. A single account is returned.
			rpcClient := rpc.Client{URL: rpcURL}
			getBalanceAndPrint(walletAccount, rpcClient)
		} else if walletIndex < 0 {
			// For when nothing is specified, shows the number of accounts in each wallet.
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
			// For when a specific wallet is specified, shows the balance of all accounts
			// in that wallet.
			checkWalletIndex()
			var accounts []string
			for address := range wallets[walletIndex].Accounts {
				accounts = append(accounts, address)
			}
			sort.Strings(accounts)
			rpcClient := rpc.Client{URL: rpcURL}
			var balanceSum, pendingSum big.Int
			for _, address := range accounts {
				balance, pending := getBalanceAndPrint(address, rpcClient)
				balanceSum.Add(&balanceSum, &balance.Int)
				pendingSum.Add(&pendingSum, &pending.Int)
			}
			if len(accounts) > 1 {
				fmt.Print(strings.Repeat(" ", 61), "Sum:")
				printAmounts(&balanceSum, &pendingSum)
			}
		}
	},
}

func getBalanceAndPrint(account string, rpcClient rpc.Client) (balance, pending *rpc.RawAmount) {
	balance, pending, err := rpcClient.AccountBalance(account)
	fatalIf(err)
	fmt.Print(account)
	printAmounts(&balance.Int, &pending.Int)
	return balance, pending
}

func printAmounts(balance, pending *big.Int) {
	if balance.Sign() > 0 {
		fmt.Printf(" %s", util.NanoAmount{Raw: balance})
	}
	if pending.Sign() > 0 {
		fmt.Printf(" (+ %s pending)", util.NanoAmount{Raw: pending})
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
