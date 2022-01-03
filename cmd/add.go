package cmd

import (
	"fmt"

	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new wallet or account",
	Run: func(cmd *cobra.Command, args []string) {
		if walletIndex < 0 {
			initNewWallet()
			fmt.Println("Added wallet.")
		} else {
			checkWalletIndex()
			wi := wallets[walletIndex]
			wi.init()
			var a *wallet.Account
			var err error
			if (walletAccountIndex > 0) {
				a, err = wi.w.NewAccount(&walletAccountIndex)
				fatalIf(err)
			} else {
				for {
					a, err = wi.w.NewAccount(nil)
					fatalIf(err)
					if _, ok := wi.Accounts[a.Address()]; !ok {
						break
					}
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
