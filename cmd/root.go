package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gonano",
	Short: "A Nano wallet written in Go",
	Long:  `Gonano is a command-line tool for managing Nano wallets.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	fatalIf(err)
}

func fatal(err ...interface{}) {
	fmt.Println(err...)
	os.Exit(1)
}

func fatalIf(err ...interface{}) {
	if err[0] != nil {
		fatal(err...)
	}
}

var rpcURL, rpcWorkURL string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gonano.yaml)")
	rootCmd.PersistentFlags().IntVarP(&walletIndex, "wallet", "w", -1, "Index of the wallet to use")
	rootCmd.PersistentFlags().StringVarP(&walletAccount, "account", "a", "", "Account to operate on")
	rootCmd.PersistentFlags().StringVarP(&rpcURL, "rpc", "r", "https://mynano.ninja/api/node", "RPC endpoint URL")
	rootCmd.PersistentFlags().StringVarP(&rpcWorkURL, "rpc-work", "s", "http://[::1]:7076", "RPC endpoint URL for work generation")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		fatalIf(err)

		// Search config in home directory with name ".gonano" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gonano")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		err = viper.SafeWriteConfig()
		fatalIf(err)
	}

	initWallets()
}
