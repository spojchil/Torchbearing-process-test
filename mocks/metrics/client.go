// Package metrics 提供实现 Metrics SDK 的 MS1 确定性内存 mock。
package metrics

import (
	"context"
	"strings"

	"github.com/spojchil/torchbearing/internal/core"
	"github.com/spojchil/torchbearing/mocks/deterministic"
	metricssdk "github.com/spojchil/torchbearing/sdk/metrics"
)

// Client 根据构造时指定的场景返回固定指标数据或强类型错误。
type Client struct {
	scenario deterministic.Scenario
}

// NewClient 创建不连接 Prometheus 或其他数据源的 Metrics mock。
func NewClient(scenario deterministic.Scenario) *Client {
	return &Client{scenario: scenario}
}

var _ metricssdk.Client = (*Client)(nil)

// Query 校验强类型查询后返回完全固定的数据序列。
func (c *Client) Query(ctx context.Context, query core.MetricQuery) (core.MetricResult, error) {
	if err := ctx.Err(); err != nil {
		return core.MetricResult{}, core.WrapAppError(
			core.ErrorCodeMetricsUnavailable,
			"metrics request canceled",
			true,
			err,
		)
	}
	if err := validateQuery(query); err != nil {
		return core.MetricResult{}, err
	}

	switch c.scenario {
	case deterministic.ScenarioSuccess:
		return successResult(), nil
	case deterministic.ScenarioEmpty:
		return core.MetricResult{State: core.DataStateEmpty, Series: []core.Series{}}, nil
	case deterministic.ScenarioMetricsFailure:
		return core.MetricResult{}, core.NewAppError(
			core.ErrorCodeMetricsUnavailable,
			"deterministic metrics failure",
			true,
		)
	case deterministic.ScenarioBoundary:
		return boundaryResult(), nil
	case deterministic.ScenarioAgentFailure:
		return core.MetricResult{}, core.NewAppError(
			core.ErrorCodeInternal,
			"metrics mock must not run after agent failure",
			false,
		)
	default:
		return core.MetricResult{}, core.NewAppError(
			core.ErrorCodeInternal,
			"unsupported deterministic metrics scenario",
			false,
		)
	}
}

func validateQuery(query core.MetricQuery) error {
	if strings.TrimSpace(query.DatasourceUID) == "" {
		return core.NewAppError(core.ErrorCodeInvalidArgument, "metrics datasource UID is required", false)
	}
	if strings.TrimSpace(query.PromQL) == "" {
		return core.NewAppError(core.ErrorCodeInvalidArgument, "PromQL is required", false)
	}
	if strings.TrimSpace(query.TimeRange.From) == "" || strings.TrimSpace(query.TimeRange.To) == "" {
		return core.NewAppError(core.ErrorCodeInvalidArgument, "metrics time range is required", false)
	}
	return nil
}

func successResult() core.MetricResult {
	return core.MetricResult{
		State: core.DataStatePresent,
		Series: []core.Series{
			{
				Name: "checkout_request_rate",
				Labels: []core.Label{
					{Name: "service", Value: "checkout"},
					{Name: "instance", Value: "checkout-01"},
				},
				Points: []core.DataPoint{
					{Timestamp: "2026-07-13T00:00:00Z", Value: 120},
					{Timestamp: "2026-07-13T00:05:00Z", Value: 128},
					{Timestamp: "2026-07-13T00:10:00Z", Value: 124},
				},
			},
		},
	}
}

func boundaryResult() core.MetricResult {
	return core.MetricResult{
		State: core.DataStatePresent,
		Series: []core.Series{
			{
				Name:   "checkout_request_total",
				Labels: []core.Label{{Name: "service", Value: "checkout"}},
				Points: []core.DataPoint{{Timestamp: "2026-07-13T00:00:00Z", Value: 0}},
			},
		},
	}
}
