package wallet

import (
	"errors"
	"sync"

	"github.com/coming-chat/wallet-SDK/core/aptos"
	"github.com/coming-chat/wallet-SDK/core/sui"
)

// Deprecated: 这个钱包对象缓存了助记词、密码、私钥等信息，继续使用有泄露资产的风险 ⚠️
type Wallet struct {
	Mnemonic string

	Keystore string
	password string

	// cache
	multiAccounts sync.Map
	aptosAccount  *aptos.Account
	suiAccount    *sui.Account
	WatchAddress  string
}

func NewWalletWithMnemonic(mnemonic string) (*Wallet, error) {
	if !IsValidMnemonic(mnemonic) {
		return nil, ErrInvalidMnemonic
	}
	return &Wallet{Mnemonic: mnemonic}, nil
}

func WatchWallet(address string) (*Wallet, error) {
	chainType := ChainTypeFrom(address)
	if chainType.Count() == 0 {
		return nil, errors.New("Invalid wallet address")
	}
	return &Wallet{WatchAddress: address}, nil
}

func (w *Wallet) IsMnemonicWallet() bool {
	return len(w.Mnemonic) > 0
}

func (w *Wallet) IsKeystoreWallet() bool {
	return len(w.Keystore) > 0
}

func (w *Wallet) IsWatchWallet() bool {
	return len(w.WatchAddress) > 0
}

func (w *Wallet) GetWatchWallet() *WatchAccount {
	return &WatchAccount{address: w.WatchAddress}
}

// Get or create the aptos account.
func (w *Wallet) GetOrCreateAptosAccount() (*aptos.Account, error) {
	cache := w.aptosAccount
	if cache != nil {
		return cache, nil
	}
	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := aptos.NewAccountWithMnemonic(w.Mnemonic)
	if err != nil {
		return nil, err
	}
	// save to cache
	w.aptosAccount = account
	return account, nil
}

// Get or create the sui account.
func (w *Wallet) GetOrCreateSuiAccount() (*sui.Account, error) {
	cache := w.suiAccount
	if cache != nil {
		return cache, nil
	}
	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := sui.NewAccountWithMnemonic(w.Mnemonic)
	if err != nil {
		return nil, err
	}
	// save to cache
	w.suiAccount = account
	return account, nil
}
