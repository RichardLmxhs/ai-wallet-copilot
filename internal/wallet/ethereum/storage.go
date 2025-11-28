package ethereum

//func (e *EthereumChain) saveBalanceToDB(ctx context.Context, address string, wei *big.Int) {
//	db := postgres.GlobalDB
//	lastBlock, err := e.LatestBlock(ctx)
//	if err != nil {
//		return
//	}
//
//	ETHPrice, err := binance.GetPriceNow(ctx, []string{CoinPriceToUSDT})
//	if err != nil {
//		logger.Global().WithContext(ctx).Error("GetPriceNow error", zap.Error(err))
//		return
//	}
//
//	totalPrice, err := CalcTotalValue(ETHPrice[CoinPriceToUSDT], wei)
//	if err != nil {
//		logger.Global().WithContext(ctx).Error("CalcTotalValue error", zap.Error(err))
//		return
//	}
//
//	t, _ := BigFloatToDecimal(totalPrice)
//
//	data := models.Wallet{
//		Address:       address,
//		Chain:         ChainName,
//		LastBlock:     &lastBlock,
//		TotalValueUSD: t,
//		TokenCount:    *wei,
//		LastIndexedAt: time.Now(),
//		Metadata:      nil,
//		CreatedAt:     time.Now(),
//		UpdatedAt:     time.Now(),
//	}
//	err = db.Clauses(clause.OnConflict{
//		Columns: []clause.Column{{Name: "address"}},
//		DoUpdates: clause.AssignmentColumns([]string{
//			"last_block",
//			"total_value_usd",
//			"token_count",
//			"nft_count",
//			"last_indexed_at",
//			"metadata",
//			"updated_at",
//		})}).Create(&data).Error
//	if err != nil {
//		logger.Global().WithContext(ctx).Error("db update error", zap.Error(err))
//		return
//	}
//
//}
//
//// LatestBlock 获取最新区块高度
//func (e *EthereumChain) LatestBlock(ctx context.Context) (int64, error) {
//	header, err := e.Client.HeaderByNumber(ctx, nil)
//	if err != nil {
//		logger.Global().WithContext(ctx).Error("get latest block error", zap.Error(err))
//		return 0, err
//	}
//	return header.Number.Int64(), nil
//}
