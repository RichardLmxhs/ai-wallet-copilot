package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/RichardLmxhs/ai-wallet-copilot/pkg/logger"
	"go.uber.org/zap"
)

type Wallet struct {
	Client      *http.Client
	MaxRetries  int
	BackoffBase time.Duration
	Endpoint    string
	APIKey      string
}

func NewWallet() *Wallet {
	return &Wallet{
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
		MaxRetries:  3,
		BackoffBase: 10 * time.Second,
		Endpoint:    config.GlobalCfg.Alchemy.Endpoint,
		APIKey:      config.GlobalCfg.Alchemy.APIKey,
	}
}

// CallRPC calls a method with params and unmarshals result into resultPtr (pointer).
// params can be slice or object (will be used as JSON params array or object depending on method).
func (w *Wallet) CallRPC(ctx context.Context, id uint64, method string, params any, resultPtr any) error {
	reqBody := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
	bs, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal rpc request: %w", err)
	}

	// prepare http request
	fullURL := path.Join(config.GlobalCfg.Alchemy.RPCURL, config.GlobalCfg.Alchemy.APIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bs))
	if err != nil {
		return fmt.Errorf("new http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// retry loop
	var lastErr error
	for attempt := 0; attempt <= w.MaxRetries; attempt++ {
		// respect context deadline
		resp, err := w.Client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http do (attempt %d): %w", attempt, err)
			// backoff if retryable
			if attempt < w.MaxRetries {
				sleep := w.BackoffBase * (1 << attempt)
				select {
				case <-time.After(sleep):
				case <-ctx.Done():
					return ctx.Err()
				}
				continue
			}
			return lastErr
		}

		// read body
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("read body: %w", err)
			if attempt < w.MaxRetries {
				sleep := w.BackoffBase * (1 << attempt)
				select {
				case <-time.After(sleep):
				case <-ctx.Done():
					return ctx.Err()
				}
				continue
			}
			return lastErr
		}

		// status not 200
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("non-200 status: %d body: %s", resp.StatusCode, string(body))
			// consider some 5xx as retryable
			if resp.StatusCode >= 500 && attempt < w.MaxRetries {
				sleep := w.BackoffBase * (1 << attempt)
				select {
				case <-time.After(sleep):
				case <-ctx.Done():
					return ctx.Err()
				}
				continue
			}
			return lastErr
		}

		// parse rpc response
		var rpcResp jsonRPCResponse
		if err := json.Unmarshal(body, &rpcResp); err != nil {
			return fmt.Errorf("unmarshal rpc response: %w (body:%s)", err, string(body))
		}
		if rpcResp.Error != nil {
			// RPC returned error
			return rpcResp.Error
		}
		if rpcResp.Result == nil {
			return errors.New("empty rpc result")
		}

		// unmarshal result into resultPtr
		if resultPtr != nil {
			if err := json.Unmarshal(*rpcResp.Result, resultPtr); err != nil {
				return fmt.Errorf("unmarshal result into ptr: %w", err)
			}
		}
		return nil
	}

	if lastErr == nil {
		lastErr = errors.New("unknown rpc error")
	}
	return lastErr
}

func (w *Wallet) GetWalletDetail(ctx context.Context, address string, networks []string) (*WalletDetail, error) {
	// 先从缓存获取
	walletDetail, err := w.GetWalletDetailCache(ctx, address)
	if err == nil {
		return walletDetail, nil
	}

	// 从alchemy获取token余额
	tokenResp, err := w.GetWalletToken(ctx, address, networks)
	if err != nil {
		return nil, err
	}

	// 从alchemy获取nft余额
	nftResp, err := w.GetWalletNFT(ctx, address, networks)
	if err != nil {
		return nil, err
	}

	// 构建返回结构
	wallet, err := w.BuildWalletDetail(ctx, address, tokenResp, nftResp)

	err = w.SetWalletDetail(ctx, address, wallet)
	if err != nil {
		return nil, err
	}

	return wallet, err
}

// GetWalletToken 从alchemy获取钱包token余额
func (w *Wallet) GetWalletToken(ctx context.Context, address string, networks []string) (*WalletTokensBalanceResponse, error) {
	chainUrl, err := url.JoinPath(w.Endpoint, fmt.Sprintf("/data/v1/%s/assets/tokens/by-address", w.APIKey))
	if err != nil {
		logger.Global().WithContext(ctx).Error("chainUrl join path error", zap.Error(err))
		return nil, err
	}

	payload := &WalletBalanceRequest{
		Addresses: []Addresses{{address, networks}},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Global().WithContext(ctx).Error("json error", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chainUrl, bytes.NewBuffer(jsonPayload))
	resp, err := w.Client.Do(req)
	if err != nil {
		logger.Global().WithContext(ctx).Error("request chain for query wallet token err", zap.Error(err))
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logger.Global().WithContext(ctx).Error("request chain for query wallet token return !200", zap.String("body", string(body)))
		return nil, errors.New("request status is not 200")
	}

	result := &WalletTokensBalanceResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetWalletNFT 从alchemy获取钱包nft余额
func (w *Wallet) GetWalletNFT(ctx context.Context, address string, networks []string) (*WalletNFTResponse, error) {
	chainUrl, err := url.JoinPath(w.Endpoint, fmt.Sprintf("/data/v1/%s/assets/nfts/by-address", w.APIKey))
	if err != nil {
		logger.Global().WithContext(ctx).Error("chainUrl join path error", zap.Error(err))
		return nil, err
	}

	payload := &WalletBalanceRequest{
		Addresses:    []Addresses{{address, networks}},
		WithMetaData: false,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Global().WithContext(ctx).Error("json error", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chainUrl, bytes.NewBuffer(jsonPayload))
	resp, err := w.Client.Do(req)
	if err != nil {
		logger.Global().WithContext(ctx).Error("request chain for query wallet nft err", zap.Error(err))
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logger.Global().WithContext(ctx).Error("request chain for query wallet nft return !200", zap.String("body", string(body)))
		return nil, errors.New("request status is not 200")
	}

	result := &WalletNFTResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetWalletTrans 从alchemy获取钱包历史交易信息
func (w *Wallet) GetWalletTrans(ctx context.Context, address string) (*[]WalletTransfersResponse, error) {
	resp := &[]WalletTransfersResponse{}
	// 转出+转入
	param := []WalletTransfersRequest{
		{
			WithMetadata: true,
			Category:     category,
			FromAddress:  address,
			ToAddress:    "",
		},
		{
			WithMetadata: true,
			Category:     category,
			FromAddress:  "",
			ToAddress:    address,
		},
	}
	if err := w.CallRPC(ctx, 2, "alchemy_getAssetTransfers", param, resp); err != nil {
		logger.Global().Error("GetWalletTrans CallContext error", zap.Error(err))
		return nil, err
	}
	return resp, nil
}

func (w *Wallet) BuildWalletDetail(
	ctx context.Context,
	address string,
	tokenResp *WalletTokensBalanceResponse,
	nftResp *WalletNFTResponse) (*WalletDetail, error) {
	detail := &WalletDetail{
		UserAddress: address,
		ChainData:   map[string]ChainData{},
	}
	totalValue := &big.Float{}

	// 构造代币信息
	for _, token := range tokenResp.Data.Tokens {
		tokenPrice := &big.Float{}
		for _, p := range token.TokenPrices {
			if p.Currency == "usd" { // 只储存美元兑换比例
				f, _, err := big.ParseFloat(p.Value, 10, 0, big.ToNearestEven)
				if err != nil {
					logger.Global().WithContext(ctx).Error("can not parse float string from token request", zap.Error(err))
					return nil, err
				}
				tokenPrice = f
				totalValue = totalValue.Add(totalValue, f)
			}
		}
		chainData, ok := detail.ChainData[token.Network]
		if !ok {
			chainData = ChainData{
				Tokens:        make([]TokenDetail, 0),
				NFTs:          make([]NFTDetail, 0),
				NFTTotalCount: 0,
			}
		}
		if token.Address == nil { // 当token address空时，说明是原生代币，例如ETH
			chainData.NativeToken = &NativeToken{
				TokenBalance: token.TokenBalance,
				TokenPrices:  tokenPrice,
			}
		} else { // 非原生代币，例如USDC
			tempToken := TokenDetail{
				TokenAddress:  token.TokenAddress,
				TokenBalance:  token.TokenBalance,
				TokenMetadata: token.TokenMetadata,
				TokenPrices:   tokenPrice,
			}
			chainData.Tokens = append(chainData.Tokens, tempToken)
		}
		detail.ChainData[token.Network] = chainData
	}
	detail.TotalValue = totalValue

	// 构造NFT信息
	nftTotalCount := 0
	for _, nft := range nftResp.Data.OwnedNfts {
		nftTotalCount++
		chainData, ok := detail.ChainData[nft.Network]
		if !ok {
			chainData = ChainData{
				Tokens:        make([]TokenDetail, 0),
				NFTs:          make([]NFTDetail, 0),
				NFTTotalCount: 0,
			}
		}
		chainData.NFTs = append(chainData.NFTs, NFTDetail{
			ContractAddress: nft.ContractAddress,
			TokenID:         nft.TokenID,
			Balance:         nft.Balance,
			Network:         nft.Network,
			Address:         nft.Address,
		})
		chainData.NFTTotalCount++
		detail.ChainData[nft.Network] = chainData
	}
	return detail, nil
}
