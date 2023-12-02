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
			getBalanceAndPrint(walletAccount)
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
			var balanceSum, receivableSum big.Int
			for _, address := range accounts {
				balance, receivable := getBalanceAndPrint(address)
				balanceSum.Add(&balanceSum, &balance.Int)
				receivableSum.Add(&receivableSum, &receivable.Int)
			}
			if len(accounts) > 1 {
				fmt.Print(strings.Repeat(" ", 61), "Sum:")
				printAmounts(&balanceSum, &receivableSum)
			}
		}
	},
}

func getBalanceAndPrint(account string) (balance, receivable *rpc.RawAmount) {
	rpcClient := rpc.Client{URL: rpcURL}
	balance, receivable, err := rpcClient.AccountBalance(account)
	fatalIf(err)
	fmt.Print(account)
	printAmounts(&balance.Int, &receivable.Int)
	return balance, receivable
}

func printAmounts(balance, receivable *big.Int) {
	if balance.Sign() > 0 {
		fmt.Printf(" %s", util.NanoAmount{Raw: balance})
	}
	if receivable.Sign() > 0 {
		fmt.Printf(" (+ %s receivable)", util.NanoAmount{Raw: receivable})
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
