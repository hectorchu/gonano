package cmd

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ssh/terminal"
)

// walletCmd represents the wallet command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Wallet management",
	Long:  `Commands for managing wallets.`,
}

type walletInfo struct {
	w        *wallet.Wallet
	Seed     string
	Accounts []string
}

var wallets []*walletInfo
var walletIndex int
var walletAccount string

func init() {
	rootCmd.AddCommand(walletCmd)
	walletCmd.PersistentFlags().IntVarP(&walletIndex, "wallet-index", "i", -1, "Index of the wallet to use")
}

func initWallets() {
	v := viper.GetStringMap("wallets")
	wallets = make([]*walletInfo, len(v))
	for i := range wallets {
		wallets[i] = &walletInfo{
			Seed:     viper.GetString(fmt.Sprintf("wallets.%d.seed", i)),
			Accounts: viper.GetStringSlice(fmt.Sprintf("wallets.%d.accounts", i)),
		}
		sort.Strings(wallets[i].Accounts)
	}
}

func checkWalletIndex() {
	if walletIndex < 0 || walletIndex >= len(wallets) {
		fatal("wallet index out of range")
	}
}

func (wi *walletInfo) init() {
	newBip39Wallet := func() {
		fmt.Print("Enter passphrase: ")
		password, err := terminal.ReadPassword(0)
		fmt.Println()
		fatalIf(err)
		wi.w, err = wallet.NewBip39Wallet(wi.Seed, string(password))
		fatalIf(err)
	}
	if wi.Seed == "" {
		entropy, err := bip39.NewEntropy(256)
		fatalIf(err)
		wi.Seed, err = bip39.NewMnemonic(entropy)
		fatalIf(err)
		newBip39Wallet()
		fmt.Println("Your secret words are:", wi.Seed)
	} else {
		seed, err := hex.DecodeString(wi.Seed)
		if err == nil {
			if len(seed) != 32 {
				fatal("invalid seed length")
			}
			wi.w, err = wallet.NewWallet(seed)
			fatalIf(err)
		} else {
			newBip39Wallet()
		}
	}
	wi.Accounts = nil
	for _, account := range wi.w.GetAccounts() {
		wi.Accounts = append(wi.Accounts, account.Address())
	}
	sort.Strings(wi.Accounts)
	fmt.Printf("%d account(s) found.\n", len(wi.Accounts))
	for i := range wallets {
		if wi == wallets[i] {
			viper.Set(fmt.Sprintf("wallets.%d", i), wi)
			err := viper.WriteConfig()
			fatalIf(err)
		}
	}
	return
}
