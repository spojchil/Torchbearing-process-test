package agent

import (
	"context"
	"reflect"
	"testing"
	"time"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
	"github.com/spojchil/torchbearing/pkg/mcp/contract"
)

type recordingTools struct {
	calls []string
}

func (t *recordingTools) SearchMetrics(context.Context, domain.ActorContext, contract.SearchMetricsInput) ([]domain.MetricCandidate, error) {
	t.calls = append(t.calls, contract.ToolMetricsSearch)
	return []domain.MetricCandidate{{Name: "http_requests_total"}}, nil
}

func (t *recordingTools) DescribeMetric(context.Context, domain.ActorContext, contract.DescribeMetricInput) (domain.MetricDescriptor, error) {
	t.calls = append(t.calls, contract.ToolMetricsDescribe)
	return domain.MetricDescriptor{Name: "http_requests_total", Type: "counter"}, nil
}

func (t *recordingTools) QueryRange(context.Context, domain.ActorContext, domain.QueryRangeRequest) (domain.QueryResult, error) {
	t.calls = append(t.calls, contract.ToolPrometheusQueryRange)
	return domain.QueryResult{Status: domain.QueryStatusSuccess, Duration: 10 * time.Millisecond, SeriesCount: 1}, nil
}

func TestRunnerUsesAllThreeMCPToolsInOrder(t *testing.T) {
	tools := &recordingTools{}
	runner := NewRunner(tools)
	request := domain.AnalysisRequest{
		Text: "checkout request rate",
		Scope: domain.AnalysisScope{
			DatasourceUID: "prometheus-main",
			TimeRange:     domain.TimeRange{From: "now-30m", To: "now"},
		},
	}
	actor := domain.ActorContext{Access: domain.AccessScope{AllowedDatasourceUIDs: []string{"prometheus-main"}}}

	result, err := runner.Analyze(context.Background(), actor, request)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	wantCalls := []string{contract.ToolMetricsSearch, contract.ToolMetricsDescribe, contract.ToolPrometheusQueryRange}
	if !reflect.DeepEqual(tools.calls, wantCalls) {
		t.Fatalf("tool calls = %v, want %v", tools.calls, wantCalls)
	}
	wantPromQL := `sum(rate(http_requests_total{service="checkout"}[5m]))`
	if result.Chart.PromQL != wantPromQL || result.Query.SeriesCount != 1 || len(result.Evidence.Metrics) != 1 {
		t.Fatalf("incomplete analysis result: %#v", result)
	}
}
