package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new wallet or account",
	Run: func(cmd *cobra.Command, args []string) {
		if walletIndex < 0 {
			seed := readPassword("Enter seed or bip39 mnemonic (leave blank for random): ")
			wi := &walletInfo{Seed: string(seed)}
			wallets = append(wallets, wi)
			wi.initNew()
			fmt.Println("Added wallet.")
		} else {
			checkWalletIndex()
			wi := wallets[walletIndex]
			wi.init()
			a, err := wi.w.NewAccount(nil)
			fatalIf(err)
			wi.Accounts[a.Address()] = a.Index()
			wi.save()
			fmt.Println("Added account", a.Address())
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
