package ethereum

import (
	"context"
	"math/big"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/chain"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/redis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumChain struct {
	*ethclient.Client
}

func NewChain() (*EthereumChain, error) {
	url := config.GlobalCfg.Chains.Ethereum.RPCURL
	rpcClient, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthereumChain{rpcClient}, err
}

// GetBalance 获取ETH余额(带缓存和数据库)
func (e *EthereumChain) GetBalance(ctx context.Context, address string, forceUpdate bool) (*big.Int, error) {
	cacheKey := ETHCachePrefix + address

	// 1. 尝试从Redis缓存读取
	if cachedBalance, err := redis.GlobalRDB.Get(ctx, cacheKey).Result(); err == nil && !forceUpdate {
		balance := new(big.Int)
		balance.SetString(cachedBalance, 10)
		return balance, nil
	}

	// 2. 从链上获取最新数据
	addr := common.HexToAddress(address)
	wei, err := e.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, err
	}

	// 3. 更新Redis缓存
	redis.GlobalRDB.Set(ctx, cacheKey, wei.String(), BalanceCacheTTL)

	// 4. 更新PostgreSQL
	go e.saveBalanceToDB(ctx, address, wei)

	return wei, nil
}

var _ chain.Chain = EthereumChain{}
