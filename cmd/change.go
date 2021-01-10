package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

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
		hash, err := a.ChangeRep(context.TODO(), args[0])
		fatalIf(err)
		fmt.Println(strings.ToUpper(hex.EncodeToString(hash)))
	},
}

func init() {
	rootCmd.AddCommand(changeCmd)
}
