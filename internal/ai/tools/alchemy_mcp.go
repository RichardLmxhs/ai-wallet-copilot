package tools

import (
	"context"
	"fmt"
	"log"

	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

var GlobalAlchemyClient *AlchemyMCPClient

type AlchemyMCPClient struct {
	client *client.Client
}

func GetAlchemyMCPTool(ctx context.Context) []tool.BaseTool {
	alchemyMcp := GlobalAlchemyClient
	err := alchemyMcp.client.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}

	_, err = alchemyMcp.client.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatal(err)
	}

	tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: alchemyMcp.client})
	if err != nil {
		log.Fatal(err)
	}
	return tools
}

// NewAlchemyMCPClient 创建一个新的 Alchemy MCP 客户端
func NewAlchemyMCPClient(apiKey string) (*AlchemyMCPClient, error) {
	// 设置环境变量
	env := []string{fmt.Sprintf("ALCHEMY_API_KEY=%s", apiKey)}

	// 创建 STDIO 客户端
	// 使用 npx 启动 Alchemy MCP 服务器
	mcpClient, err := client.NewStdioMCPClient(
		"npx",
		env,
		"-y",
		"@alchemy/mcp-server",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}

	// 初始化连接
	ctx := context.Background()
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Alchemy Go Client",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	if _, err := mcpClient.Initialize(ctx, initRequest); err != nil {
		mcpClient.Close()
		return nil, fmt.Errorf("failed to initialize client: %w", err)
	}
	c := &AlchemyMCPClient{
		client: mcpClient,
	}

	GlobalAlchemyClient = c

	return c, nil
}
