package ai

// 定义钱包分析的提示模板
const WalletAnalysisPrompt = `
目标：请你作为 Web3/安全专家，分析一个钱包的资产与行为。

【基本信息】
- 总资产：{{.TotalAssetValue}}
- Token：
  {{range .Tokens}}
  - {{.Symbol}}: {{.Value}}
  {{end}}
- NFT：{{.NFTCount}}个

【历史交易】
{{.TransactionSummary}}

【风险标记】
{{.RiskMarkers}}

请提供以下分析内容:
1. 钱包行为总结
2. 风险评分 (0-100)
3. 风险分析
4. 操作建议
`
