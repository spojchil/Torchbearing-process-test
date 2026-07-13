package chart

import (
	"errors"
	"testing"

	"github.com/spojchil/torchbearing/internal/core"
)

// TestBuilderScenarios 覆盖成功、空结果、状态矛盾失败和单点边界图表。
func TestBuilderScenarios(t *testing.T) {
	tests := []struct {
		name       string
		plan       core.AnalysisPlan
		result     core.MetricResult
		wantCharts int
		wantCode   core.ErrorCode
	}{
		{name: "success", plan: validPlan(core.ChartTypeTimeSeries), result: presentResult(3), wantCharts: 1},
		{name: "empty", plan: validPlan(core.ChartTypeTimeSeries), result: core.MetricResult{State: core.DataStateEmpty, Series: []core.Series{}}, wantCharts: 0},
		{name: "failure", plan: validPlan(core.ChartTypeTimeSeries), result: core.MetricResult{State: core.DataStatePresent, Series: []core.Series{}}, wantCode: core.ErrorCodeNoData},
		{name: "boundary", plan: validPlan(core.ChartTypeStat), result: presentResult(1), wantCharts: 1},
	}

	builder := NewBuilder()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			charts, err := builder.Build(test.plan, test.result)
			if test.wantCode != "" {
				assertBuilderErrorCode(t, err, test.wantCode)
				if charts != nil {
					t.Fatalf("charts = %+v, want nil", charts)
				}
				return
			}

			if err != nil {
				t.Fatalf("Build returned unexpected error: %v", err)
			}
			if len(charts) != test.wantCharts {
				t.Fatalf("chart count = %d, want %d", len(charts), test.wantCharts)
			}
			if charts == nil {
				t.Fatal("charts must be a non-nil deterministic slice")
			}
			if len(charts) == 1 && charts[0].Type != test.plan.Chart.Type {
				t.Fatalf("chart type = %q, want %q", charts[0].Type, test.plan.Chart.Type)
			}
		})
	}
}

// TestBuilderRejectsUnknownChartType 验证未知图表类型不会越过稳定契约。
func TestBuilderRejectsUnknownChartType(t *testing.T) {
	plan := validPlan(core.ChartType("unknown"))

	_, err := NewBuilder().Build(plan, presentResult(1))
	assertBuilderErrorCode(t, err, core.ErrorCodeInvalidArgument)
}

func validPlan(chartType core.ChartType) core.AnalysisPlan {
	return core.AnalysisPlan{
		Message: "已生成 checkout 服务请求速率图表。",
		Query: core.MetricQuery{
			DatasourceUID: "prometheus-mock",
			PromQL:        `sum(rate(http_requests_total{service="checkout"}[5m]))`,
			TimeRange:     core.TimeRange{From: "now-30m", To: "now"},
		},
		Chart: core.ChartHint{ID: "checkout-request-rate", Title: "Checkout request rate", Type: chartType},
	}
}

func presentResult(pointCount int) core.MetricResult {
	points := make([]core.DataPoint, pointCount)
	return core.MetricResult{
		State: core.DataStatePresent,
		Series: []core.Series{
			{Name: "checkout", Labels: []core.Label{}, Points: points},
		},
	}
}

func assertBuilderErrorCode(t *testing.T, err error, want core.ErrorCode) {
	t.Helper()

	var typed *core.AppError
	if !errors.As(err, &typed) {
		t.Fatalf("error type = %T, want *core.AppError", err)
	}
	if typed.Code != want {
		t.Fatalf("error code = %q, want %q", typed.Code, want)
	}
}
