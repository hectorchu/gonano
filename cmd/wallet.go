package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/hectorchu/gonano/wallet"
	"github.com/spf13/viper"
	"github.com/tyler-smith/go-bip39"
)

type walletInfo struct {
	w          *wallet.Wallet
	Seed, Salt string
	IsBip39    bool
	Accounts   map[string]uint32
}

var wallets []*walletInfo
var walletIndex int
var walletAccount string

func initWallets() {
	v := viper.GetStringMap("wallets")
	wallets = make([]*walletInfo, len(v))
	for i := range wallets {
		key := func(s string) string {
			return fmt.Sprintf("wallets.%d.%s", i, s)
		}
		wallets[i] = &walletInfo{
			Seed:     viper.GetString(key("seed")),
			Salt:     viper.GetString(key("salt")),
			IsBip39:  viper.GetBool(key("isbip39")),
			Accounts: make(map[string]uint32),
		}
		for k, v := range viper.GetStringMap(key("accounts")) {
			wallets[i].Accounts[k] = uint32(v.(int))
		}
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
	wi.w.RPC.URL = rpcURL
	wi.w.RPCWork.URL = rpcWorkURL
	err := wi.w.ScanForAccounts()
	fatalIf(err)
	for _, a := range wi.w.GetAccounts() {
		wi.Accounts[a.Address()] = a.Index()
	}
	wi.save()
}

func (wi *walletInfo) save() {
	for i := range wallets {
		if wi == wallets[i] {
			viper.Set(fmt.Sprintf("wallets.%d", i), wi)
			err := viper.WriteConfig()
			fatalIf(err)
		}
	}
}

func getAccount() (a *wallet.Account) {
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
	return
}
