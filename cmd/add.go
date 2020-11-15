package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new wallet or account",
}

func init() {
	rootCmd.AddCommand(addCmd)
}
