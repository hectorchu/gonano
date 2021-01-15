package cmd

import (
	"fmt"

	"github.com/hectorchu/gonano/util"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an amount of Nano from an account",
	Long: `Send an amount of Nano from an account.

  send <destination> <amount>`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		a := getAccount()
		amount, err := util.NanoAmountFromString(args[1])
		fatalIf(err)
		hash, err := a.Send(args[0], amount.Raw)
		fatalIf(err)
		fmt.Println(hash)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
