package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/models"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/storage/postgres"
	"github.com/RichardLmxhs/ai-wallet-copilot/pkg/logger"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

// StoreWalletDetail 存储钱包详情到数据库
func (w *Wallet) StoreWalletDetail(ctx context.Context, walletAddress string, detail *WalletDetail) error {
	// 参数验证
	if detail == nil {
		logger.Global().WithContext(ctx).Error("wallet detail is nil")
		return fmt.Errorf("wallet detail is nil")
	}

	if walletAddress == "" {
		logger.Global().WithContext(ctx).Error("wallet address is empty")
		return fmt.Errorf("wallet address is empty")
	}

	if detail.ChainData == nil {
		logger.Global().WithContext(ctx).Error("wallet chain data is nil")
		return fmt.Errorf("wallet chain data is nil")
	}

	db := postgres.GlobalDB
	wallets := make([]models.Wallet, 0)
	walletTokens := make([]models.WalletToken, 0)
	walletNFTs := make([]models.WalletNFT, 0)
	priceCaches := make([]models.PriceCache, 0)

	// 准备当前时间戳，确保所有记录使用相同的更新时间
	now := time.Now()

	// 准备数据源标识
	source := "alchemy"

	for network, chainData := range detail.ChainData {
		// 处理wallet数据
		chainCoin, err := strconv.ParseInt(chainData.NativeToken.TokenBalance, 10, 64)
		if err != nil {
			logger.Global().WithContext(ctx).Error("convert TokenBalance error", zap.Error(err), zap.String("chain", network))
			continue // 跳过当前链，继续处理下一个链
		}

		// 创建钱包元数据
		walletMetadata, _ := json.Marshal(map[string]interface{}{
			"total_value_usd": detail.TotalValue.String(),
			"token_count":     len(chainData.Tokens),
			"nft_count":       len(chainData.NFTs),
		})

		w := models.Wallet{
			Address:       walletAddress,
			Chain:         network,
			ChainCoins:    chainCoin,
			LastIndexedAt: now,
			Metadata:      walletMetadata,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		wallets = append(wallets, w)

		// 处理原生代币价格 - 存储到price_cache
		if chainData.NativeToken.TokenPrices != nil {
			// 将big.Float转换为decimal.Decimal
			priceStr := chainData.NativeToken.TokenPrices.String()
			priceUSD, err := decimal.NewFromString(priceStr)
			if err == nil {
				// 原生代币的token_address为nil
				priceCache := models.PriceCache{
					TokenAddress: nil, // 原生代币没有地址
					Chain:        network,
					PriceUSD:     &priceUSD,
					Source:       &source,
					UpdatedAt:    now,
				}
				priceCaches = append(priceCaches, priceCache)
			}
		}

		// 处理wallet_tokens数据和token价格
		for _, token := range chainData.Tokens {
			// 解析余额为decimal.Decimal
			balance, err := decimal.NewFromString(token.TokenBalance)
			if err != nil {
				logger.Global().WithContext(ctx).Error("parse token balance error",
					zap.Error(err),
					zap.String("token_address", token.TokenAddress))
				continue // 跳过当前token，继续处理下一个
			}

			// 处理小数位数
			decimals := 18 // 默认值
			if token.TokenMetadata.Decimals != nil {
				decimals = *token.TokenMetadata.Decimals
			}

			// 创建token元数据
			tokenMetadata, _ := json.Marshal(token.TokenMetadata)

			tokenModel := models.WalletToken{
				WalletAddress: walletAddress,
				Chain:         network,
				TokenAddress:  token.TokenAddress,
				TokenType:     "ERC20", // 默认类型
				Balance:       balance,
				Decimals:      decimals,
				Symbol:        token.TokenMetadata.Symbol,
				Name:          token.TokenMetadata.Name,
				Metadata:      tokenMetadata,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			walletTokens = append(walletTokens, tokenModel)

			// 处理token价格 - 存储到price_cache
			if token.TokenPrices != nil {
				// 将big.Float转换为decimal.Decimal
				priceStr := token.TokenPrices.String()
				priceUSD, err := decimal.NewFromString(priceStr)
				if err == nil {
					tokenAddress := token.TokenAddress
					priceCache := models.PriceCache{
						TokenAddress: &tokenAddress,
						Chain:        network,
						PriceUSD:     &priceUSD,
						Source:       &source,
						Symbol:       token.TokenMetadata.Symbol,
						Logo:         token.TokenMetadata.Logo,
						UpdatedAt:    now,
					}
					priceCaches = append(priceCaches, priceCache)
				}
			}
		}

		// 处理wallet_nft数据
		for _, nft := range chainData.NFTs {
			// 创建NFT元数据
			nftMetadata, _ := json.Marshal(map[string]interface{}{
				"balance": nft.Balance,
			})

			nftModel := models.WalletNFT{
				WalletAddress:   walletAddress,
				Chain:           network,
				ContractAddress: nft.ContractAddress,
				TokenID:         nft.TokenID,
				TokenType:       "ERC721", // 默认类型
				Metadata:        nftMetadata,
				CreatedAt:       now,
				UpdatedAt:       now,
			}
			walletNFTs = append(walletNFTs, nftModel)
		}
	}

	// 使用事务批量处理
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.Global().WithContext(ctx).Error("panic during database operation", zap.Any("recover", r))
		}
	}()

	// 存储wallet数据
	for _, wallet := range wallets {
		// 使用Upsert操作，当地址存在时更新，不存在时插入
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "address"}, {Name: "chain"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"chain_coins",
				"last_indexed_at",
				"metadata",
				"updated_at",
			}),
		}).Create(&wallet).Error; err != nil {
			tx.Rollback()
			logger.Global().WithContext(ctx).Error("save wallet error", zap.Error(err))
			return err
		}
	}

	// 存储wallet_tokens数据
	for _, token := range walletTokens {
		// 使用Upsert操作，基于wallet_address和token_address的联合唯一索引
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "wallet_address"}, {Name: "token_address"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"balance",
				"decimals",
				"symbol",
				"name",
				"metadata",
				"updated_at",
			}),
		}).Create(&token).Error; err != nil {
			tx.Rollback()
			logger.Global().WithContext(ctx).Error("save wallet token error",
				zap.Error(err),
				zap.String("token_address", token.TokenAddress))
			return err
		}
	}

	// 存储wallet_nft数据
	for _, nft := range walletNFTs {
		// 使用Upsert操作，基于wallet_address、contract_address和token_id的组合
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "wallet_address"}, {Name: "contract_address"}, {Name: "token_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"token_type",
				"metadata",
				"updated_at",
			}),
		}).Create(&nft).Error; err != nil {
			tx.Rollback()
			logger.Global().WithContext(ctx).Error("save wallet nft error",
				zap.Error(err),
				zap.String("contract_address", nft.ContractAddress),
				zap.String("token_id", nft.TokenID))
			return err
		}
	}

	// 存储price_cache数据
	for _, priceCache := range priceCaches {
		// 使用Upsert操作，基于token_address和chain的组合
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "token_address"}, {Name: "chain"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"price_usd",
				"source",
				"symbol",
				"logo",
				"updated_at",
				"metadata",
			}),
		}).Create(&priceCache).Error; err != nil {
			tx.Rollback()
			tokenAddr := "native"
			if priceCache.TokenAddress != nil {
				tokenAddr = *priceCache.TokenAddress
			}
			logger.Global().WithContext(ctx).Error("save price cache error",
				zap.Error(err),
				zap.String("token_address", tokenAddr),
				zap.String("chain", priceCache.Chain))
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Global().WithContext(ctx).Error("commit transaction error", zap.Error(err))
		return err
	}

	logger.Global().WithContext(ctx).Info("wallet detail stored successfully",
		zap.String("wallet_address", walletAddress),
		zap.Int("wallet_count", len(wallets)),
		zap.Int("token_count", len(walletTokens)),
		zap.Int("nft_count", len(walletNFTs)),
		zap.Int("price_cache_count", len(priceCaches)))

	return nil
}
