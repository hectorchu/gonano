package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a wallet",
	Long:  `Create a wallet with a supplied or random seed.`,
	Run: func(cmd *cobra.Command, args []string) {
		seed := readPassword("Enter seed or bip39 mnemonic (leave blank for random): ")
		wi := &walletInfo{Seed: string(seed)}
		wallets = append(wallets, wi)
		wi.initNew()
		fmt.Println("Added wallet.")
	},
}

func init() {
	walletCmd.AddCommand(createCmd)
}
