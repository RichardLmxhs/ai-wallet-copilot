package wallet

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/redis"
)

const (
	WalletPrefix               = "wallet:detail:"
	WalletNativeTokenTTLPrefix = "wallet:NativeToken:TTL:%s"
	WalletTokenTTLPrefix       = "wallet:Token:TTL:{%s}"
	WalletNFTTTLPrefix         = "wallet:NFT:TTL:{%s}"

	DefaultWalletTTL = 30 * time.Minute
	DefaultETHTTL    = 1 * time.Minute
	DefaultTokenTTL  = 5 * time.Minute
	DefaultNFTTTL    = 10 * time.Minute
)

type WalletDetail struct {
	UserAddress string               `json:"userAddress"`
	TotalValue  *big.Float           `json:"totalValue"`
	ChainData   map[string]ChainData `json:"chainData"`
}

type ChainData struct {
	NativeToken   *NativeToken  `json:"nativeToken"`
	Tokens        []TokenDetail `json:"tokens"` // {network:[]tokens}
	NFTs          []NFTDetail   `json:"nfts"`   // {network: []nfts}
	NFTTotalCount int           `json:"nftTotalCount"`
}

type NativeToken struct {
	TokenBalance string     `json:"tokenBalance"`
	TokenPrices  *big.Float `json:"tokenPrices"`
}

type TokenDetail struct {
	TokenAddress  string        `json:"tokenAddress"`
	TokenBalance  string        `json:"tokenBalance"`
	TokenMetadata TokenMetadata `json:"tokenMetadata"`
	TokenPrices   *big.Float    `json:"tokenPrices"`
}

type NFTDetail struct {
	ContractAddress string `json:"contractAddress"`
	TokenID         string `json:"tokenId"`
	Balance         string `json:"balance"`
	Network         string `json:"network"`
	Address         string `json:"address"`
}

func (w *Wallet) GetWalletDetailCache(ctx context.Context, walletAddress string) (*WalletDetail, error) {
	detail := &WalletDetail{}
	err := redis.GlobalRDB.Get(ctx, fmt.Sprintf(WalletPrefix+walletAddress)).Scan(detail)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

// SetWalletDetail 存储钱包缓存到redis
func (w *Wallet) SetWalletDetail(ctx context.Context, walletAddress string, detail *WalletDetail) error {
	err := redis.GlobalRDB.Set(ctx, fmt.Sprintf(WalletPrefix+walletAddress), detail, DefaultWalletTTL).Err()
	if err != nil {
		return err
	}

	err = redis.GlobalRDB.Set(ctx, fmt.Sprintf(WalletNativeTokenTTLPrefix+walletAddress),
		time.Now().Format(time.RFC3339), -1).Err()
	if err != nil {
		return err
	}

	err = redis.GlobalRDB.Set(ctx, fmt.Sprintf(WalletTokenTTLPrefix+walletAddress),
		time.Now().Format(time.RFC3339), -1).Err()
	if err != nil {
		return err
	}

	err = redis.GlobalRDB.Set(ctx, fmt.Sprintf(WalletNFTTTLPrefix+walletAddress),
		time.Now().Format(time.RFC3339), -1).Err()
	if err != nil {
		return err
	}
	return nil
}
