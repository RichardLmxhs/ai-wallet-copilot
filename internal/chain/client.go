package chain

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/RichardLmxhs/ai-wallet-copilot/pkg/logger"
	"go.uber.org/zap"
)

type Chain interface {
}

type WalletTokenRequest struct {
	Addresses []Addresses `json:"addresses"`
}

type Addresses struct {
	Address  string   `json:"address"`
	Networks []string `json:"networks"`
}

type WalletTokenResponse struct {
	Data Data `json:"data"`
}
type TokenMetadata struct {
	Decimals int    `json:"decimals"`
	Logo     string `json:"logo"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
}
type TokenPrices struct {
	Currency      string    `json:"currency"`
	Value         string    `json:"value"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}
type Tokens struct {
	Address       string        `json:"address"`
	Network       string        `json:"network"`
	TokenAddress  string        `json:"tokenAddress"`
	TokenBalance  string        `json:"tokenBalance"`
	TokenMetadata TokenMetadata `json:"tokenMetadata"`
	TokenPrices   []TokenPrices `json:"tokenPrices"`
	Error         string        `json:"error"`
}
type Data struct {
	Tokens  []Tokens `json:"tokens"`
	PageKey string   `json:"pageKey"`
}

func GetWalletTokenBalance(ctx context.Context, address string, forceUpdate bool, networks []string) (*WalletTokenResponse, error) {
	endpoint := config.GlobalCfg.Alchemy.Endpoint
	apiKey := config.GlobalCfg.Alchemy.APIKey
	chainUrl, err := url.JoinPath(endpoint, fmt.Sprintf("/data/v1/%s/assets/tokens/by-address", apiKey))
	if err != nil {
		logger.Global().WithContext(ctx).Error("chainUrl join path error", zap.Error(err))
		return nil, err
	}
	client := http.Client{
		Timeout: 15 * time.Second,
	}

	payload := &WalletTokenRequest{
		Addresses: []Addresses{{address, networks}},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Global().WithContext(ctx).Error("json error", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chainUrl, bytes.NewBuffer(jsonPayload))
	resp, err := client.Do(req)
	if err != nil {
		logger.Global().WithContext(ctx).Error("request chain for query wallet token err", zap.Error(err))
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logger.Global().WithContext(ctx).Error("request chain for query wallet token return !200", zap.String("body", string(body)))
		return nil, errors.New("request status is not 200")
	}

	result := &WalletTokenResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
