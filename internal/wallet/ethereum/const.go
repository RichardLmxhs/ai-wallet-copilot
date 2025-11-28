package ethereum

import "time"

const (
	ChainName       = "ethereum"
	CoinPriceToUSDT = "ETHUSDT"

	ETHCachePrefix = "eth:balance:"

	BalanceCacheTTL   = 30 * time.Second // 余额缓存30秒
	ERC20CacheTTL     = 5 * time.Minute  // ERC20缓存5分钟
	NFTCacheTTL       = 10 * time.Minute // NFT缓存10分钟
	TxHistoryCacheTTL = 2 * time.Minute  // 交易历史缓存2分钟
	ContractLabelTTL  = 24 * time.Hour   // 合约标签缓存24小时
)
