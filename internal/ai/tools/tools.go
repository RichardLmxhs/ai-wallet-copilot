package tools

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/mcp"

	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
)

func getMCPTool(ctx context.Context) []tool.BaseTool {
	alchemyMcp, err := NewAlchemyMCPClient("xxx")
	if err != nil {
		log.Fatal(err)
	}
	err = alchemyMcp.client.Start(ctx)
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
