package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "Change representative for an account",
	Long: `Change representative for an account.

  change <address>`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		a := getAccount()
		hash, err := a.ChangeRep(args[0])
		fatalIf(err)
		fmt.Println(hash)
	},
}

func init() {
	rootCmd.AddCommand(changeCmd)
}
