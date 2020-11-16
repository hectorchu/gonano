package cmd

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an amount of Nano from an account",
	Long: `Send an amount of Nano from an account.

send <destination> <amount>`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		checkWalletAccount()
		if walletIndex < 0 {
			for i, wi := range wallets {
				if _, ok := wi.Accounts[walletAccount]; ok {
					walletIndex = i
					break
				}
			}
			if walletIndex < 0 {
				fatal("account not found in any wallet")
			}
		}
		checkWalletIndex()
		wi := wallets[walletIndex]
		wi.init()
		index, ok := wi.Accounts[walletAccount]
		if !ok {
			fatal("account not found in the specified wallet")
		}
		a, err := wi.w.NewAccount(&index)
		fatalIf(err)
		amount, err := wallet.NanoAmountFromString(args[1])
		fatalIf(err)
		hash, err := a.Send(args[0], amount.Raw)
		fatalIf(err)
		fmt.Println(strings.ToUpper(hex.EncodeToString(hash)))
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
