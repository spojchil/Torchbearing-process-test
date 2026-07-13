// Package agent 定义 MS1 的 Agent SDK 边界，具体 Agent 适配器由其他模块负责实现。
package agent

import (
	"context"

	"github.com/spojchil/torchbearing/internal/core"
)

// Client 将带分析范围的自然语言请求转换为强类型分析计划。
type Client interface {
	Plan(ctx context.Context, request core.AgentRequest) (core.AnalysisPlan, error)
}
