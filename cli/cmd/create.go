package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		wi.init()
		viper.Set(fmt.Sprintf("wallets.%d", len(wallets)), wi)
		wallets = append(wallets, wi)
		err = viper.WriteConfig()
		fatalIf(err)
		fmt.Println("Added wallet.")
	},
}

func init() {
	walletCmd.AddCommand(createCmd)
}
