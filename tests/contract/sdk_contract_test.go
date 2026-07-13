package contract_test

import (
	"context"
	"testing"

	"github.com/spojchil/torchbearing/internal/core"
	agentsdk "github.com/spojchil/torchbearing/sdk/agent"
	metricssdk "github.com/spojchil/torchbearing/sdk/metrics"
)

// fixedAgentClient 是仅用于契约测试的确定性 Agent fake。
type fixedAgentClient struct {
	plan core.AnalysisPlan
	err  error
}

func (c fixedAgentClient) Plan(context.Context, core.AgentRequest) (core.AnalysisPlan, error) {
	return c.plan, c.err
}

// fixedMetricsClient 是仅用于契约测试的确定性 Metrics fake。
type fixedMetricsClient struct {
	result core.MetricResult
	err    error
}

func (c fixedMetricsClient) Query(context.Context, core.MetricQuery) (core.MetricResult, error) {
	return c.result, c.err
}

// 编译期断言确保测试 fake 始终满足 A 定义的 SDK 接口。
var _ agentsdk.Client = fixedAgentClient{}
var _ metricssdk.Client = fixedMetricsClient{}

// TestSDKContractsCarryTypedSuccessValues 验证 SDK 可以传递完整强类型成功结果。
func TestSDKContractsCarryTypedSuccessValues(t *testing.T) {
	query := core.MetricQuery{
		DatasourceUID: "prometheus-mock",
		PromQL:        `sum(rate(http_requests_total{service="checkout"}[5m]))`,
		TimeRange:     core.TimeRange{From: "now-30m", To: "now"},
	}
	agent := fixedAgentClient{plan: core.AnalysisPlan{Query: query}}
	metrics := fixedMetricsClient{result: core.MetricResult{State: core.DataStateEmpty}}

	plan, err := agent.Plan(context.Background(), core.AgentRequest{})
	if err != nil {
		t.Fatalf("Plan returned unexpected error: %v", err)
	}
	if plan.Query.PromQL != query.PromQL {
		t.Fatalf("PromQL = %q, want %q", plan.Query.PromQL, query.PromQL)
	}

	result, err := metrics.Query(context.Background(), plan.Query)
	if err != nil {
		t.Fatalf("Query returned unexpected error: %v", err)
	}
	if result.State != core.DataStateEmpty {
		t.Fatalf("State = %q, want %q", result.State, core.DataStateEmpty)
	}
}

// TestSDKContractsCarryTypedFailures 验证 SDK 可以原样传递强类型失败。
func TestSDKContractsCarryTypedFailures(t *testing.T) {
	agent := fixedAgentClient{err: core.NewAppError(
		core.ErrorCodeAgentUnavailable,
		"deterministic agent failure",
		true,
	)}

	_, err := agent.Plan(context.Background(), core.AgentRequest{})
	typed, ok := err.(*core.AppError)
	if !ok {
		t.Fatalf("error type = %T, want *core.AppError", err)
	}
	if typed.Code != core.ErrorCodeAgentUnavailable || !typed.Retryable {
		t.Fatalf("typed error = %+v, want retryable AGENT_UNAVAILABLE", typed)
	}
}
