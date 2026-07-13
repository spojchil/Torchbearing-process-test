package agent

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/spojchil/torchbearing/internal/core"
	"github.com/spojchil/torchbearing/mocks/deterministic"
)

// TestClientPlanScenarios 覆盖成功、空结果准备、Agent 失败和单值边界计划。
func TestClientPlanScenarios(t *testing.T) {
	tests := []struct {
		name          string
		scenario      deterministic.Scenario
		wantCode      core.ErrorCode
		wantChartType core.ChartType
		wantQueryPart string
	}{
		{name: "success", scenario: deterministic.ScenarioSuccess, wantChartType: core.ChartTypeTimeSeries, wantQueryPart: `service="checkout"`},
		{name: "empty", scenario: deterministic.ScenarioEmpty, wantChartType: core.ChartTypeTimeSeries, wantQueryPart: `service="unknown-service"`},
		{name: "failure", scenario: deterministic.ScenarioAgentFailure, wantCode: core.ErrorCodeAgentUnavailable},
		{name: "boundary", scenario: deterministic.ScenarioBoundary, wantChartType: core.ChartTypeStat, wantQueryPart: "sum(http_requests_total"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			plan, err := NewClient(test.scenario).Plan(context.Background(), validRequest())
			if test.wantCode != "" {
				assertErrorCode(t, err, test.wantCode)
				if plan != (core.AnalysisPlan{}) {
					t.Fatalf("plan = %+v, want zero value", plan)
				}
				return
			}

			if err != nil {
				t.Fatalf("Plan returned unexpected error: %v", err)
			}
			if plan.Chart.Type != test.wantChartType {
				t.Fatalf("chart type = %q, want %q", plan.Chart.Type, test.wantChartType)
			}
			if !strings.Contains(plan.Query.PromQL, test.wantQueryPart) {
				t.Fatalf("PromQL = %q, want substring %q", plan.Query.PromQL, test.wantQueryPart)
			}
			if plan.Query.DatasourceUID != "prometheus-mock" {
				t.Fatalf("datasource UID = %q, want prometheus-mock", plan.Query.DatasourceUID)
			}
		})
	}
}

// TestClientRejectsEmptyText 验证空自然语言输入不会生成伪造计划。
func TestClientRejectsEmptyText(t *testing.T) {
	request := validRequest()
	request.Text = "  "

	_, err := NewClient(deterministic.ScenarioSuccess).Plan(context.Background(), request)
	assertErrorCode(t, err, core.ErrorCodeInvalidArgument)
}

func validRequest() core.AgentRequest {
	return core.AgentRequest{
		Text: "查看 checkout 服务过去 30 分钟的请求速率",
		Context: core.AnalysisContext{Scope: core.AnalysisScope{
			DatasourceUID: "prometheus-mock",
			TimeRange:     core.TimeRange{From: "now-30m", To: "now"},
		}},
	}
}

func assertErrorCode(t *testing.T, err error, want core.ErrorCode) {
	t.Helper()

	var typed *core.AppError
	if !errors.As(err, &typed) {
		t.Fatalf("error type = %T, want *core.AppError", err)
	}
	if typed.Code != want {
		t.Fatalf("error code = %q, want %q", typed.Code, want)
	}
}
