package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/wallet"
)

// 定义钱包分析的提示模板
const SystemRolePrompt = `
你是一个 加密钱包安全分析专家、安全风控工程师、链上行为分析师。
你的职责是基于用户提供的钱包结构化数据，进行：

资产画像分析：分链分析资产规模、分布、价值来源

交易行为分析：识别活跃度、资产流向、潜在资金来源与风险

风险监控：识别风险地址交互、高风险代币、大额异常转账、洗币特征

综合安全评分：0–100 分

给出明确结论：钱包是否安全、是否异常、是否被盗风险

给出可执行建议（如：是否需转移资产、是否需取消授权、是否需停止交互）

你的输出必须基于用户提供的数据，不进行虚构数据补充。

你在输出内容时：

注重结构化、清晰、完整

分段输出：资产概览 / 行为分析 / 风险识别 / 风险等级 / 建议

尽量量化信息（如风险评分、单链占比、交易频率）

不需要描述数据结构和代码，只分析链上信息

绝不生成不基于输入推导的内容

如果用户缺少某些数据，你需要指出数据缺失情况，并说明无法分析的部分。
`

func BuildWalletUserPrompt(walletDetail *wallet.WalletDetail) string {
	var sb strings.Builder

	// 基础信息
	sb.WriteString("以下是用户钱包的链上数据，请基于这些结构化信息进行资产分析、风险识别和安全建议。\n\n")

	sb.WriteString("### 🧾 WalletDetail\n")
	sb.WriteString(fmt.Sprintf("- UserAddress: %s\n", walletDetail.UserAddress))
	if walletDetail.TotalValue != nil {
		sb.WriteString(fmt.Sprintf("- TotalValue: %s\n", walletDetail.TotalValue.Text('f', 6)))
	}

	// 链数据
	sb.WriteString("\n### 🏦 Chain Data\n")

	for chain, data := range walletDetail.ChainData {
		sb.WriteString(fmt.Sprintf("\n#### Chain: %s\n\n", chain))

		// Native Token
		if data.NativeToken != nil {
			sb.WriteString("• Native Token\n")
			sb.WriteString(fmt.Sprintf("  - Balance: %s\n", data.NativeToken.TokenBalance))
			if data.NativeToken.TokenPrices != nil {
				sb.WriteString(fmt.Sprintf("  - Price: %s\n", data.NativeToken.TokenPrices.Text('f', 6)))
			}
		}

		// Tokens
		sb.WriteString("\n• Tokens:\n")
		for _, t := range data.Tokens {
			meta, _ := json.Marshal(t.TokenMetadata)
			sb.WriteString(fmt.Sprintf("  - TokenAddress: %s\n", t.TokenAddress))
			sb.WriteString(fmt.Sprintf("    Balance: %s\n", t.TokenBalance))
			sb.WriteString(fmt.Sprintf("    Metadata: %s\n", string(meta)))
			if t.TokenPrices != nil {
				sb.WriteString(fmt.Sprintf("    Price: %s\n", t.TokenPrices.Text('f', 6)))
			}
		}

		// NFTs
		sb.WriteString("\n• NFTs:\n")
		for _, nft := range data.NFTs {
			sb.WriteString(fmt.Sprintf("  - Contract: %s\n", nft.ContractAddress))
			sb.WriteString(fmt.Sprintf("    TokenID: %s\n", nft.TokenID))
			sb.WriteString(fmt.Sprintf("    Balance: %s\n", nft.Balance))
			sb.WriteString(fmt.Sprintf("    Network: %s\n", nft.Network))
			sb.WriteString(fmt.Sprintf("    Address: %s\n", nft.Address))
		}

		sb.WriteString(fmt.Sprintf("\n• NFT Total Count: %d\n", data.NFTTotalCount))
	}

	// Transfers
	sb.WriteString("\n### 🔁 Transfers (可选)\n")
	if walletDetail.Transfers != nil {
		js, _ := json.MarshalIndent(walletDetail.Transfers, "", "  ")
		sb.WriteString(string(js) + "\n")
	} else {
		sb.WriteString("无\n")
	}

	// 分析任务
	sb.WriteString(`
---

请分析：

1. 资产分布与构成  
2. 链间资产占比  
3. 代币风险（honeypot、空气币、假稳定币、垃圾代币）  
4. 交易行为分析：活跃度、大额转账、可疑流向  
5. 授权风险（若可推断）  
6. 综合风险评分（0–100）  
7. 用户应采取的操作建议（如撤销授权、迁移资产等）

请只基于以上输入数据，不要虚构不存在的数据。
`)

	return sb.String()
}
