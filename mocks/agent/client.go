// Package agent 提供实现 Agent SDK 的 MS1 确定性 mock。
package agent

import (
	"context"
	"strings"

	"github.com/spojchil/torchbearing/internal/core"
	"github.com/spojchil/torchbearing/mocks/deterministic"
	agentsdk "github.com/spojchil/torchbearing/sdk/agent"
)

// Client 根据构造时指定的场景返回固定分析计划或强类型错误。
type Client struct {
	scenario deterministic.Scenario
}

// NewClient 创建不访问真实大模型或第三方 API 的 Agent mock。
func NewClient(scenario deterministic.Scenario) *Client {
	return &Client{scenario: scenario}
}

var _ agentsdk.Client = (*Client)(nil)

// Plan 将已验证的分析上下文转换为固定 PromQL 和图表意图。
func (c *Client) Plan(ctx context.Context, request core.AgentRequest) (core.AnalysisPlan, error) {
	if err := ctx.Err(); err != nil {
		return core.AnalysisPlan{}, core.WrapAppError(
			core.ErrorCodeAgentUnavailable,
			"agent request canceled",
			true,
			err,
		)
	}
	if strings.TrimSpace(request.Text) == "" {
		return core.AnalysisPlan{}, core.NewAppError(
			core.ErrorCodeInvalidArgument,
			"analysis text is required",
			false,
		)
	}
	if err := validateScope(request.Context.Scope); err != nil {
		return core.AnalysisPlan{}, err
	}

	switch c.scenario {
	case deterministic.ScenarioSuccess, deterministic.ScenarioMetricsFailure:
		return requestRatePlan(request.Context.Scope, "checkout"), nil
	case deterministic.ScenarioEmpty:
		plan := requestRatePlan(request.Context.Scope, "unknown-service")
		plan.Message = "指定范围内没有可展示的数据。"
		return plan, nil
	case deterministic.ScenarioAgentFailure:
		return core.AnalysisPlan{}, core.NewAppError(
			core.ErrorCodeAgentUnavailable,
			"deterministic agent failure",
			true,
		)
	case deterministic.ScenarioBoundary:
		return singleValuePlan(request.Context.Scope), nil
	default:
		return core.AnalysisPlan{}, core.NewAppError(
			core.ErrorCodeInternal,
			"unsupported deterministic agent scenario",
			false,
		)
	}
}

func validateScope(scope core.AnalysisScope) error {
	if scope.DatasourceUID == "" || scope.TimeRange.From == "" || scope.TimeRange.To == "" {
		return core.NewAppError(core.ErrorCodeInvalidScope, "resolved analysis scope is required", false)
	}
	return nil
}

func requestRatePlan(scope core.AnalysisScope, service string) core.AnalysisPlan {
	return core.AnalysisPlan{
		Message: "已生成 checkout 服务请求速率图表。",
		Query: core.MetricQuery{
			DatasourceUID: scope.DatasourceUID,
			PromQL:        `sum(rate(http_requests_total{service="` + service + `"}[5m]))`,
			TimeRange:     scope.TimeRange,
		},
		Chart: core.ChartHint{
			ID:    service + "-request-rate",
			Title: "Checkout request rate",
			Type:  core.ChartTypeTimeSeries,
		},
	}
}

func singleValuePlan(scope core.AnalysisScope) core.AnalysisPlan {
	return core.AnalysisPlan{
		Message: "已生成 checkout 服务当前请求总量。",
		Query: core.MetricQuery{
			DatasourceUID: scope.DatasourceUID,
			PromQL:        `sum(http_requests_total{service="checkout"})`,
			TimeRange:     scope.TimeRange,
		},
		Chart: core.ChartHint{
			ID:    "checkout-request-total",
			Title: "Checkout request total",
			Type:  core.ChartTypeStat,
		},
	}
}
