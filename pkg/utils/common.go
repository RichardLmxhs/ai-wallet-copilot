package utils

import (
	"regexp"
	"strings"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	tron "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/mr-tron/base58"
)

const (
	ETH = "eth"
	BTC = "btc"
	TRX = "trx"
	SOL = "sol"
)

var NetworkMap = map[string][]string{
	ETH: {"eth-mainnet"},
	BTC: {"bitcoin-mainnet"},
	TRX: {"tron-mainnet"},
	SOL: {"solana-mainnet"},
}

func ValidateAddress(chain string, addr string) bool {
	switch chain {
	case ETH:
		return IsValidEthAddress(addr)
	case BTC:
		return IsValidBTCAddress(addr)
	case TRX:
		return IsValidTrxAddress(addr)
	case SOL:
		return IsValidSolAddress(addr)
	default:
		return false
	}
}

func IsValidSolAddress(addr string) bool {
	b, err := base58.Decode(addr)
	if err != nil {
		return false
	}
	return len(b) == 32
}

func IsValidTrxAddress(addr string) bool {
	_, err1 := tron.Base58ToAddress(addr)
	_, err2 := tron.Base64ToAddress(addr)
	if err1 != nil || err2 != nil {
		return true
	}
	return false
}

func IsValidBTCAddress(addr string) bool {
	_, err := btcutil.DecodeAddress(addr, nil) // 主网用 nil 即可
	return err == nil
}

func IsValidEthAddress(addr string) bool {
	if !strings.HasPrefix(addr, "0x") {
		return false
	}
	if len(addr) != 42 {
		return false
	}
	// hex 格式校验
	matched, _ := regexp.MatchString(`^0x[0-9a-fA-F]{40}$`, addr)
	if !matched {
		return false
	}
	// 使用 geth 的校验函数（包含 checksum）
	return common.IsHexAddress(addr)
}
