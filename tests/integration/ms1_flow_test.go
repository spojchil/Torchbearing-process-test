package integration_test

import (
	"context"
	"testing"

	"github.com/spojchil/torchbearing/internal/bootstrap"
	"github.com/spojchil/torchbearing/internal/contracts"
	"github.com/spojchil/torchbearing/internal/core"
	"github.com/spojchil/torchbearing/mocks/deterministic"
)

// TestMS1MainFlow 覆盖主流程成功、空结果、两类失败和单点边界。
func TestMS1MainFlow(t *testing.T) {
	tests := []struct {
		name       string
		scenario   deterministic.Scenario
		request    contracts.AnalysisRequest
		requestID  string
		chartCount int
		chartType  core.ChartType
		errorCode  core.ErrorCode
	}{
		{name: "success", scenario: deterministic.ScenarioSuccess, request: standardRequest(), requestID: "mock-analysis-001", chartCount: 1, chartType: core.ChartTypeTimeSeries},
		{name: "empty", scenario: deterministic.ScenarioEmpty, request: standardRequest(), requestID: "mock-analysis-002", chartCount: 0},
		{name: "agent failure", scenario: deterministic.ScenarioAgentFailure, request: standardRequest(), requestID: "mock-analysis-003", errorCode: core.ErrorCodeAgentUnavailable},
		{name: "metrics failure", scenario: deterministic.ScenarioMetricsFailure, request: standardRequest(), requestID: "mock-analysis-004", errorCode: core.ErrorCodeMetricsUnavailable},
		{name: "single point boundary", scenario: deterministic.ScenarioBoundary, request: standardRequest(), requestID: "mock-analysis-005", chartCount: 1, chartType: core.ChartTypeStat},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gateway, err := bootstrap.NewMS1Gateway(test.scenario)
			if err != nil {
				t.Fatalf("bootstrap gateway: %v", err)
			}

			response, failure := gateway.Analyze(context.Background(), test.request)
			if test.errorCode != "" {
				if failure == nil {
					t.Fatal("failure = nil, want typed error response")
				}
				if failure.Code != test.errorCode {
					t.Fatalf("error code = %q, want %q", failure.Code, test.errorCode)
				}
				if failure.RequestID != test.requestID {
					t.Fatalf("error request ID = %q, want %q", failure.RequestID, test.requestID)
				}
				return
			}

			if failure != nil {
				t.Fatalf("failure = %+v, want nil", failure)
			}
			if response.RequestID != test.requestID {
				t.Fatalf("request ID = %q, want %q", response.RequestID, test.requestID)
			}
			if !response.Mock {
				t.Fatal("Mock = false, want true")
			}
			if len(response.Charts) != test.chartCount {
				t.Fatalf("chart count = %d, want %d", len(response.Charts), test.chartCount)
			}
			if test.chartCount > 0 && response.Charts[0].Type != test.chartType {
				t.Fatalf("chart type = %q, want %q", response.Charts[0].Type, test.chartType)
			}
		})
	}
}

// TestMS1InvalidScopeStopsBeforeAdapters 验证 B 的关键边界错误直接结束主流程。
func TestMS1InvalidScopeStopsBeforeAdapters(t *testing.T) {
	gateway, err := bootstrap.NewMS1Gateway(deterministic.ScenarioBoundary)
	if err != nil {
		t.Fatalf("bootstrap gateway: %v", err)
	}
	request := standardRequest()
	request.Scope.TimeRange = contracts.TimeRange{
		From: "2026-07-13T01:00:00Z",
		To:   "2026-07-13T00:00:00Z",
	}

	response, failure := gateway.Analyze(context.Background(), request)
	if response.RequestID != "" || response.Message != "" || response.Charts != nil || response.Mock {
		t.Fatalf("response = %+v, want zero value", response)
	}
	if failure == nil || failure.Code != core.ErrorCodeInvalidScope {
		t.Fatalf("failure = %+v, want INVALID_SCOPE", failure)
	}
	if failure.RequestID != "mock-analysis-005" {
		t.Fatalf("request ID = %q, want mock-analysis-005", failure.RequestID)
	}
}

func standardRequest() contracts.AnalysisRequest {
	return contracts.AnalysisRequest{
		Text: "查看 checkout 服务过去 30 分钟的请求速率",
		Scope: contracts.AnalysisScope{
			DatasourceUID: "prometheus-mock",
			TimeRange:     contracts.TimeRange{From: "now-30m", To: "now"},
		},
	}
}
