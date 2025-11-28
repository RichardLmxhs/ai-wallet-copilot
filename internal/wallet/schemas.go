package wallet

import "time"

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
