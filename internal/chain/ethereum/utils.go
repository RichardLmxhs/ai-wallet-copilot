package ethereum

import (
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"
)

func CalcTotalValue(priceStr string, balanceWei *big.Int) (*big.Float, error) {
	// 1. 解析价格字符串为 big.Float
	priceFloat, _, err := big.ParseFloat(priceStr, 10, 256, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("invalid price string: %v", err)
	}

	// 2. 将 Wei 转换为 big.Float
	balanceFloat := new(big.Float).SetInt(balanceWei)

	// 3. ETH = wei / 1e18
	weiBase := new(big.Float).SetFloat64(1e18)
	balanceETH := new(big.Float).Quo(balanceFloat, weiBase)

	// 4. 总价 = ETH * price
	total := new(big.Float).Mul(balanceETH, priceFloat)

	return total, nil
}

func BigFloatToDecimal(f *big.Float) (decimal.Decimal, error) {
	// 使用f.Text生成十进制字符串，'f'表示小数格式，-1表示不限精度
	s := f.Text('f', -1)

	return decimal.NewFromString(s)
}
