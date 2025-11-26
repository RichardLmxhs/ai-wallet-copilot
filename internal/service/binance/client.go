package binance

import (
	"context"
	"errors"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/RichardLmxhs/ai-wallet-copilot/pkg/logger"
	binance_connector "github.com/binance/binance-connector-go"
	"go.uber.org/zap"
)

var BinanceClient *binance_connector.Client

func InitBinanceService(cfg *config.Config) {
	BinanceClient = binance_connector.NewClient(cfg.BinanceAPI.AK, cfg.BinanceAPI.SK, cfg.BinanceAPI.BaseURL)
	return
}

func GetPriceNow(ctx context.Context, name []string) (map[string]string, error) {
	req := BinanceClient.NewTickerPriceService()
	req = req.Symbols(name)
	result, err := req.Do(ctx)
	if err != nil {
		logger.Global().WithContext(ctx).Error("get Price from binance error:%v", zap.Error(err))
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("can not find the coin")
	}
	res := map[string]string{}
	for _, coin := range result {
		res[coin.Symbol] = coin.Price
	}
	return res, nil
}
