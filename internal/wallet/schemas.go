package wallet

import (
	"encoding/json"
	"fmt"
	"time"
)

var (
	category = []string{
		"erc20",
		"erc1155",
		"erc721",
		"external",
		"internal",
	}
	maxCount = 500
)

type WalletBalanceRequest struct {
	Addresses    []Addresses `json:"addresses"`
	WithMetaData bool        `json:"withMetadata"`
}

type Addresses struct {
	Address  string   `json:"address"`
	Networks []string `json:"networks"`
}

type WalletTokensBalanceResponse struct {
	Data TokenData `json:"data"`
}
type TokenMetadata struct {
	Decimals *int    `json:"decimals"`
	Logo     *string `json:"logo"`
	Name     *string `json:"name"`
	Symbol   *string `json:"symbol"`
}
type TokenPrices struct {
	Currency      string    `json:"currency"`
	Value         string    `json:"value"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

type Tokens struct {
	Address       *string       `json:"address"`
	Network       string        `json:"network"`
	TokenAddress  string        `json:"tokenAddress"`
	TokenBalance  string        `json:"tokenBalance"`
	TokenMetadata TokenMetadata `json:"tokenMetadata"`
	TokenPrices   []TokenPrices `json:"tokenPrices"`
	Error         string        `json:"error"`
}
type TokenData struct {
	Tokens  []Tokens `json:"tokens"`
	PageKey string   `json:"pageKey"`
}

type WalletNFTResponse struct {
	Data NFTData `json:"data"`
}
type OwnedNfts struct {
	ContractAddress     string `json:"contractAddress"`
	TokenID             string `json:"tokenId"`
	Balance             string `json:"balance"`
	IsSpam              bool   `json:"isSpam"`
	SpamClassifications []any  `json:"spamClassifications"`
	Network             string `json:"network"`
	Address             string `json:"address"`
}
type NFTData struct {
	OwnedNfts  []OwnedNfts `json:"ownedNfts"`
	TotalCount int         `json:"totalCount"`
	PageKey    any         `json:"pageKey"`
}

type WalletTransfersRequest struct {
	Category    []string `json:"category"`
	FromAddress string   `json:"fromAddress"`
	ToAddress   string   `json:"toAddress"`
	MaxCount    string   `json:"maxCount"` // eg. 0x03
}

type RawContract struct {
	Value   string `json:"value"`
	Address string `json:"address"`
	Decimal string `json:"decimal"`
}
type Metadata struct {
	BlockTimestamp time.Time `json:"blockTimestamp"`
}

type WalletTransfersResponse struct {
	Transfers []Transfers `json:"transfers"`
	PageKey   string      `json:"pageKey"`
}

type Transfers struct {
	BlockNum        string      `json:"blockNum"`
	UniqueID        string      `json:"uniqueId"`
	Hash            string      `json:"hash"`
	From            string      `json:"from"`
	To              string      `json:"to"`
	Value           float64     `json:"value"`
	Erc721TokenID   interface{} `json:"erc721TokenId"`
	Erc1155Metadata interface{} `json:"erc1155Metadata"`
	TokenID         interface{} `json:"tokenId"`
	Asset           string      `json:"asset"`
	Category        string      `json:"category"`
	RawContract     RawContract `json:"rawContract"`
}

// JsonRPC structures
type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      uint64      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type jsonRPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      uint64           `json:"id"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *rpcError        `json:"error,omitempty"`
}
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
}
