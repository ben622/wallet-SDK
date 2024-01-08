package wallet

import (
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/coming-chat/wallet-SDK/core/aptos"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/base/inter"
	"github.com/coming-chat/wallet-SDK/core/sui"
)

const (
	// contains ethereum, bsc, chainx_eth, polygon...
	ChainTypeEthereum = "ethereum"
	ChainTypeBitcoin  = "bitcoin"
	ChainTypeCosmos   = "cosmos"
	ChainTypeSolana   = "solana"

	// contains chainx, minix, sherpax, polkadot...
	ChainTypePolka    = "polka"
	ChainTypeSignet   = "signet"
	ChainTypeDoge     = "dogecoin"
	ChainTypeTerra    = "terra"
	ChainTypeAptos    = "aptos"
	ChainTypeSui      = "sui"
	ChainTypeStarcoin = "starcoin"
	ChainTypeStarknet = "starknet"
)

// Deprecated: renamed to `ChainTypeOfWatchAddress()`.
func ChainTypeFrom(address string) *base.StringArray {
	return ChainTypeOfWatchAddress(address)
}

// Only support evm, btc, cosmos, solana now.
func ChainTypeOfWatchAddress(address string) *base.StringArray {
	res := &base.StringArray{}
	for _, ch := range []byte(address) {
		valid := (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
		if !valid {
			return res
		}
	}
	if strings.HasPrefix(address, "0x") || strings.HasPrefix(address, "0X") {
		if aptos.IsValidAddress(address) {
			res.Append(ChainTypeAptos)
		}
		if sui.IsValidAddress(address) {
			res.Append(ChainTypeSui)
		}
	} else {
		// if btc.IsValidAddress(address, btc.ChainSignet) {
		// 	res.Append(ChainTypeSignet)
		// }
		// if doge.IsValidAddress(address, doge.ChainMainnet) {
		// 	res.Append(ChainTypeDoge)
		// }
	}
	return res
}

func ChainTypeOfPrivateKey(prikey string) *base.StringArray {
	res := &base.StringArray{}
	isValid, length := inter.IsHexString(prikey)
	if isValid {
		if length == 64 {
			res.Append(ChainTypeBitcoin)
			res.Append(ChainTypeEthereum)
			res.Append(ChainTypePolka)
			res.Append(ChainTypeSignet)
			res.Append(ChainTypeDoge)
			res.Append(ChainTypeCosmos)
			res.Append(ChainTypeTerra)
			res.Append(ChainTypeAptos)
			res.Append(ChainTypeSui)
			res.Append(ChainTypeStarcoin)
		}
		if length == 128 {
			res.Append(ChainTypeSolana)
		}
	} else {
		if strings.HasPrefix(prikey, "0x") || strings.HasPrefix(prikey, "0X") {
			prikey = prikey[2:] // remove 0x prefix
		}
		if !inter.IsASCII(prikey) {
			return res
		}
		wif, err := btcutil.DecodeWIF(prikey)
		if err == nil {
			if wif.IsForNet(&chaincfg.MainNetParams) {
				res.Append(ChainTypeBitcoin)
			} else {
				res.Append(ChainTypeSignet)
			}
		}
	}
	return res
}
