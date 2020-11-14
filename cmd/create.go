package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a wallet",
	Long:  `Create a wallet with a supplied or random seed.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter seed or bip39 mnemonic (leave blank for random): ")
		seed, err := terminal.ReadPassword(0)
		fmt.Println()
		fatalIf(err)
		wi := &walletInfo{Seed: string(seed)}
		wallets = append(wallets, wi)
		wi.init()
		fmt.Println("Added wallet.")
	},
}

func init() {
	walletCmd.AddCommand(createCmd)
}
