package utils

import (
	"regexp"
)

// 以太坊地址正则表达式
var ethAddressRegex = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

// IsValidEthAddress 检查是否是有效的以太坊地址
func IsValidEthAddress(address string) bool {
	return ethAddressRegex.MatchString(address)
}

// FormatAddress 格式化地址（如省略中间部分）
func FormatAddress(address string) string {
	if len(address) < 10 {
		return address
	}
	return address[0:6] + "..." + address[len(address)-4:]
}
