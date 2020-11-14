package cmd

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send from a wallet",
	Long: `Send an amount of Nano from a wallet account.

send <destination> <amount>`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if walletIndex < 0 {
			for i, wi := range wallets {
				if sort.SearchStrings(wi.Accounts, walletAccount) < len(wi.Accounts) {
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
		account := wi.w.GetAccount(walletAccount)
		if account == nil {
			fatal("account not found in any wallet")
		}
		amount, err := wallet.NanoAmountFromString(args[1])
		fatalIf(err)
		hash, err := account.Send(args[0], amount.Raw)
		fatalIf(err)
		fmt.Println(strings.ToUpper(hex.EncodeToString(hash)))
	},
}

func init() {
	walletCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVarP(&walletAccount, "account", "a", "", "Account to send from")
	sendCmd.MarkFlagRequired("account")
}
