// Package client defines the Agent Host side of the observability MCP
// boundary. ActorContext is injected separately from model-editable inputs.
package client

import (
	"context"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
)

const (
	ToolMetricsSearch        = "metrics.search"
	ToolMetricsDescribe      = "metrics.describe"
	ToolPrometheusQueryRange = "prometheus.query_range"
)

type SearchMetricsInput struct {
	Scope domain.AnalysisScope
	Text  string
	Limit int
}

type DescribeMetricInput struct {
	Scope  domain.AnalysisScope
	Metric string
}

type Client interface {
	SearchMetrics(ctx context.Context, actor domain.ActorContext, input SearchMetricsInput) ([]domain.MetricCandidate, error)
	DescribeMetric(ctx context.Context, actor domain.ActorContext, input DescribeMetricInput) (domain.MetricDescriptor, error)
	QueryRange(ctx context.Context, actor domain.ActorContext, input domain.QueryRangeRequest) (domain.QueryResult, error)
}
