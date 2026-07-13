package metrics

import (
	"context"
	"errors"
	"testing"

	"github.com/spojchil/torchbearing/internal/core"
	"github.com/spojchil/torchbearing/mocks/deterministic"
)

// TestClientQueryScenarios 覆盖成功、空结果、可预期失败和单点边界数据。
func TestClientQueryScenarios(t *testing.T) {
	tests := []struct {
		name       string
		scenario   deterministic.Scenario
		wantState  core.DataState
		wantSeries int
		wantPoints int
		wantCode   core.ErrorCode
	}{
		{name: "success", scenario: deterministic.ScenarioSuccess, wantState: core.DataStatePresent, wantSeries: 1, wantPoints: 3},
		{name: "empty", scenario: deterministic.ScenarioEmpty, wantState: core.DataStateEmpty, wantSeries: 0},
		{name: "failure", scenario: deterministic.ScenarioMetricsFailure, wantCode: core.ErrorCodeMetricsUnavailable},
		{name: "boundary", scenario: deterministic.ScenarioBoundary, wantState: core.DataStatePresent, wantSeries: 1, wantPoints: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewClient(test.scenario).Query(context.Background(), validQuery())
			if test.wantCode != "" {
				assertMetricsErrorCode(t, err, test.wantCode)
				if result.State != "" || result.Series != nil {
					t.Fatalf("result = %+v, want zero value", result)
				}
				return
			}

			if err != nil {
				t.Fatalf("Query returned unexpected error: %v", err)
			}
			if result.State != test.wantState {
				t.Fatalf("state = %q, want %q", result.State, test.wantState)
			}
			if len(result.Series) != test.wantSeries {
				t.Fatalf("series count = %d, want %d", len(result.Series), test.wantSeries)
			}
			if result.Series == nil {
				t.Fatal("series must be a non-nil deterministic slice")
			}
			if test.wantSeries > 0 && len(result.Series[0].Points) != test.wantPoints {
				t.Fatalf("point count = %d, want %d", len(result.Series[0].Points), test.wantPoints)
			}
		})
	}
}

// TestClientRejectsIncompleteQuery 验证 adapter 不会接受缺失 PromQL 的查询。
func TestClientRejectsIncompleteQuery(t *testing.T) {
	query := validQuery()
	query.PromQL = ""

	_, err := NewClient(deterministic.ScenarioSuccess).Query(context.Background(), query)
	assertMetricsErrorCode(t, err, core.ErrorCodeInvalidArgument)
}

func validQuery() core.MetricQuery {
	return core.MetricQuery{
		DatasourceUID: "prometheus-mock",
		PromQL:        `sum(rate(http_requests_total{service="checkout"}[5m]))`,
		TimeRange:     core.TimeRange{From: "now-30m", To: "now"},
	}
}

func assertMetricsErrorCode(t *testing.T, err error, want core.ErrorCode) {
	t.Helper()

	var typed *core.AppError
	if !errors.As(err, &typed) {
		t.Fatalf("error type = %T, want *core.AppError", err)
	}
	if typed.Code != want {
		t.Fatalf("error code = %q, want %q", typed.Code, want)
	}
}
