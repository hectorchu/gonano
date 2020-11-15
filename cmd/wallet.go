package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/hectorchu/gonano/rpc"
	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/viper"
	"github.com/tyler-smith/go-bip39"
)

type walletInfo struct {
	w          *wallet.Wallet
	Seed, Salt string
	IsBip39    bool
	Accounts   []string
}

var wallets []*walletInfo
var walletIndex int
var walletAccount string
var rpcClient = rpc.Client{URL: "https://mynano.ninja/api/node"}

func initWallets() {
	v := viper.GetStringMap("wallets")
	wallets = make([]*walletInfo, len(v))
	for i := range wallets {
		fmt := func(key string) string {
			return fmt.Sprintf("wallets.%d.%s", i, key)
		}
		wallets[i] = &walletInfo{
			Seed:     viper.GetString(fmt("seed")),
			Salt:     viper.GetString(fmt("salt")),
			IsBip39:  viper.GetBool(fmt("isbip39")),
			Accounts: viper.GetStringSlice(fmt("accounts")),
		}
		sort.Strings(wallets[i].Accounts)
	}
}

func checkWalletIndex() {
	if walletIndex < 0 {
		fatal("wallet index (-w) not specified")
	} else if walletIndex >= len(wallets) {
		fatal("wallet index out of range")
	}
}

func checkWalletAccount() {
	if walletAccount == "" {
		fatal("wallet account (-a) not specified")
	}
}

func (wi *walletInfo) initNew() {
	password := readPassword("Enter password: ")
	password2 := readPassword("Re-enter password: ")
	if !bytes.Equal(password, password2) {
		fatal("password mismatch")
	}
	key, salt, err := deriveKey(password, nil)
	fatalIf(err)
	wi.Salt = hex.EncodeToString(salt)
	initBip39 := func(entropy []byte) {
		enc, err := encrypt(entropy, key)
		fatalIf(err)
		wi.Seed = hex.EncodeToString(enc)
		wi.IsBip39 = true
		wi.initBip39(entropy, password)
	}
	if wi.Seed == "" {
		entropy, err := bip39.NewEntropy(256)
		fatalIf(err)
		mnemonic, err := bip39.NewMnemonic(entropy)
		fatalIf(err)
		fmt.Println("Your secret words are:", mnemonic)
		initBip39(entropy)
	} else {
		seed, err := hex.DecodeString(wi.Seed)
		if err == nil {
			enc, err := encrypt(seed, key)
			fatalIf(err)
			wi.Seed = hex.EncodeToString(enc)
			wi.initRegularSeed(seed)
		} else {
			entropy, err := bip39.EntropyFromMnemonic(wi.Seed)
			fatalIf(err)
			initBip39(entropy)
		}
	}
	wi.initAccounts()
}

func (wi *walletInfo) init() {
	password := readPassword("Enter password: ")
	enc, err := hex.DecodeString(wi.Seed)
	fatalIf(err)
	salt, err := hex.DecodeString(wi.Salt)
	fatalIf(err)
	key, _, err := deriveKey(password, salt)
	fatalIf(err)
	seed, err := decrypt(enc, key)
	fatalIf(err)
	if wi.IsBip39 {
		wi.initBip39(seed, password)
	} else {
		wi.initRegularSeed(seed)
	}
	wi.initAccounts()
}

func (wi *walletInfo) initRegularSeed(seed []byte) {
	if len(seed) != 32 {
		fatal("invalid seed length")
	}
	var err error
	wi.w, err = wallet.NewWallet(seed)
	fatalIf(err)
}

func (wi *walletInfo) initBip39(entropy, password []byte) {
	mnemonic, err := bip39.NewMnemonic(entropy)
	fatalIf(err)
	wi.w, err = wallet.NewBip39Wallet(mnemonic, string(password))
	fatalIf(err)
}

func (wi *walletInfo) initAccounts() {
	accounts := make(map[string]bool)
	for _, account := range wi.Accounts {
		accounts[account] = true
	}
	for _, account := range wi.w.GetAccounts() {
		accounts[account.Address()] = true
	}
	wi.Accounts = nil
	for account := range accounts {
		wi.Accounts = append(wi.Accounts, account)
	}
	wi.save()
}

func (wi *walletInfo) save() {
	sort.Strings(wi.Accounts)
	for i := range wallets {
		if wi == wallets[i] {
			viper.Set(fmt.Sprintf("wallets.%d", i), wi)
			err := viper.WriteConfig()
			fatalIf(err)
		}
	}
}
