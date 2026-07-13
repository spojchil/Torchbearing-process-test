// Package agent contains the AgentRunner adapter. The deterministic Runner in
// the skeleton proves orchestration only; Issue #6 replaces its decision logic
// with an Eino-backed implementation behind the same application port.
package agent

import (
	"context"
	"fmt"
	"time"

	application "github.com/spojchil/torchbearing/pkg/application/analysis"
	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
	mcpclient "github.com/spojchil/torchbearing/pkg/mcp/client"
)

type Runner struct {
	tools mcpclient.Client
}

func NewRunner(tools mcpclient.Client) *Runner {
	return &Runner{tools: tools}
}

func (r *Runner) Analyze(ctx context.Context, actor domain.ActorContext, request domain.AnalysisRequest) (domain.AnalysisResult, error) {
	candidates, err := r.tools.SearchMetrics(ctx, actor, mcpclient.SearchMetricsInput{
		Scope: request.Scope,
		Text:  request.Text,
		Limit: 5,
	})
	if err != nil {
		return domain.AnalysisResult{}, err
	}
	if len(candidates) == 0 {
		return domain.AnalysisResult{}, domain.NewError(domain.CodeMetricNotFound, "no matching metric was found")
	}

	descriptor, err := r.tools.DescribeMetric(ctx, actor, mcpclient.DescribeMetricInput{
		Scope:  request.Scope,
		Metric: candidates[0].Name,
	})
	if err != nil {
		return domain.AnalysisResult{}, err
	}

	expression := fmt.Sprintf(`sum(rate(%s{service="checkout"}[5m]))`, descriptor.Name)
	queryResult, err := r.tools.QueryRange(ctx, actor, domain.QueryRangeRequest{
		DatasourceUID: request.Scope.DatasourceUID,
		Expression:    expression,
		TimeRange:     request.Scope.TimeRange,
		Step:          15 * time.Second,
		MaxSeries:     100,
	})
	if err != nil {
		return domain.AnalysisResult{}, err
	}
	if queryResult.Status == domain.QueryStatusNoData || queryResult.SeriesCount == 0 {
		return domain.AnalysisResult{}, domain.NewError(domain.CodeNoData, "query returned no data")
	}

	return domain.AnalysisResult{
		Chart: domain.ChartSpec{
			ID:            "request-rate",
			Title:         "Checkout 请求速率（架构 Stub）",
			Type:          domain.ChartTypeTimeseries,
			DatasourceUID: request.Scope.DatasourceUID,
			PromQL:        expression,
			TimeRange:     request.Scope.TimeRange,
			Legend:        "checkout",
		},
		Query: domain.QuerySummary{
			Language:    "promql",
			Expression:  expression,
			Status:      queryResult.Status,
			Duration:    queryResult.Duration,
			SeriesCount: queryResult.SeriesCount,
		},
		Evidence: domain.Evidence{
			Metrics:     []string{descriptor.Name},
			Explanation: "架构 Stub 已依次调用指标搜索、指标描述和范围查询三个 MCP 工具。",
		},
	}, nil
}

var _ application.AgentRunner = (*Runner)(nil)
