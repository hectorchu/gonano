package cmd

import (
	"context"
	"fmt"

	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new wallet or account",
	Run: func(cmd *cobra.Command, args []string) {
		if walletIndex < 0 {
			initNewWallet(context.TODO())
			fmt.Println("Added wallet.")
		} else {
			checkWalletIndex()

			wi := wallets[walletIndex]

			wi.init()

			var a *wallet.Account

			for {
				var err error

				a, err = wi.w.NewAccount(nil)
				fatalIf(err)

				if _, ok := wi.Accounts[a.Address()]; !ok {
					break
				}
			}

			wi.Accounts[a.Address()] = a.Index()
			wi.save()

			fmt.Println("Added account", a.Address())
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
