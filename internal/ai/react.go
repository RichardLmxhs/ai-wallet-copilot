package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	cfg "github.com/RichardLmxhs/ai-wallet-copilot/internal/config"
	"github.com/RichardLmxhs/ai-wallet-copilot/internal/wallet"
	"github.com/RichardLmxhs/ai-wallet-copilot/pkg/logger"
	clc "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/cozeloop-go"
	"go.uber.org/zap"

	"github.com/RichardLmxhs/ai-wallet-copilot/internal/ai/tools"
)

type WalletAIAnalyze struct {
	Tools         []tool.InvokableTool
	WalletDetail  *wallet.WalletDetail
	AnalyzeStatus AgentStatus
}

type AgentStatus struct {
	Status string
	Result string
}

func NewWalletAIAnalyze(toolsName []string, walletDetail *wallet.WalletDetail) *WalletAIAnalyze {
	w := &WalletAIAnalyze{
		AnalyzeStatus: AgentStatus{
			Status: StatusStart,
			Result: "",
		},
	}
	for _, toolName := range toolsName {
		if t, ok := tools.ToolsMap[toolName]; ok {
			w.Tools = append(w.Tools, t())
		}
	}
	w.WalletDetail = walletDetail
	return w
}

func (w *WalletAIAnalyze) Run(ctx context.Context) {
	agentTools := []tool.BaseTool{}
	for _, t := range w.Tools {
		agentTools = append(agentTools, t)
	}

	cozeloopApiToken := cfg.GlobalCfg.AI.CozeloopAPIToken
	cozeloopWorkspaceID := cfg.GlobalCfg.AI.CozeWorkSpaceID
	arkApiKey := cfg.GlobalCfg.AI.ARKAPIKey
	arkModelName := cfg.GlobalCfg.AI.ARKModelName

	var handlers []callbacks.Handler
	if cozeloopApiToken != "" && cozeloopWorkspaceID != "" {
		client, err := cozeloop.NewClient(
			cozeloop.WithAPIToken(cozeloopApiToken),
			cozeloop.WithWorkspaceID(cozeloopWorkspaceID),
		)
		if err != nil {
			panic(err)
		}
		defer client.Close(ctx)
		handlers = append(handlers, clc.NewLoopHandler(client))
	}
	callbacks.AppendGlobalHandlers(handlers...)

	config := &ark.ChatModelConfig{
		APIKey: arkApiKey,
		Model:  arkModelName,
	}

	arkModel, err := ark.NewChatModel(ctx, config)
	if err != nil {
		logger.Global().WithContext(ctx).Error("failed to create chat model:", zap.Error(err))
		return
	}

	ragent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: arkModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: agentTools,
		},
		// StreamToolCallChecker: toolCallChecker, // uncomment it to replace the default tool call checker with custom one
	})
	if err != nil {
		logger.Global().WithContext(ctx).Error("failed to create agent: ", zap.Error(err))
		return
	}
	opt := []agent.AgentOption{
		agent.WithComposeOptions(compose.WithCallbacks(&LoggerCallback{})),
		//react.WithChatModelOptions(ark.WithCache(cacheOption)),
	}

	sr, err := ragent.Stream(ctx, []*schema.Message{
		{
			Role:    schema.System,
			Content: SystemRolePrompt,
		},
		{
			Role:    schema.User,
			Content: BuildWalletUserPrompt(w.WalletDetail),
		},
	}, opt...)
	if err != nil {
		logger.Global().WithContext(ctx).Error("failed to stream: %v", zap.Error(err))
		return
	}

	defer sr.Close() // remember to close the stream

	logger.Global().WithContext(ctx).Info("\n\n===== start streaming =====\n\n")

	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			// error
			logger.Global().WithContext(ctx).Warn("failed to recv: %v", zap.Error(err))
			return
		}

		// 打字机打印
		logs.Tokenf("%v", msg.Content)
	}

	logs.Infof("\n\n===== finished =====\n")
	time.Sleep(2 * time.Second)
}

type LoggerCallback struct {
	callbacks.HandlerBuilder // 可以用 callbacks.HandlerBuilder 来辅助实现 callback
}

func (cb *LoggerCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	fmt.Println("==================")
	inputStr, _ := json.MarshalIndent(input, "", "  ")
	fmt.Printf("[OnStart] %s\n", string(inputStr))
	return ctx
}

func (cb *LoggerCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	fmt.Println("=========[OnEnd]=========")
	outputStr, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(outputStr))
	return ctx
}

func (cb *LoggerCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Println("=========[OnError]=========")
	fmt.Println(err)
	return ctx
}

func (cb *LoggerCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {

	var graphInfoName = react.GraphName

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("[OnEndStream] panic err:", err)
			}
		}()

		defer output.Close() // remember to close the stream in defer

		fmt.Println("=========[OnEndStream]=========")
		for {
			frame, err := output.Recv()
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			if err != nil {
				fmt.Printf("internal error: %s\n", err)
				return
			}

			s, err := json.Marshal(frame)
			if err != nil {
				fmt.Printf("internal error: %s\n", err)
				return
			}

			if info.Name == graphInfoName { // 仅打印 graph 的输出, 否则每个 stream 节点的输出都会打印一遍
				fmt.Printf("%s: %s\n", info.Name, string(s))
			}
		}

	}()
	return ctx
}

func (cb *LoggerCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo,
	input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	defer input.Close()
	return ctx
}
