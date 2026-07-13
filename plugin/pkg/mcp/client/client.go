// Package client defines the Agent Host side of the observability MCP
// boundary. ActorContext is injected separately from model-editable inputs.
package client

import (
	"context"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
	"github.com/spojchil/torchbearing/pkg/mcp/contract"
)

type Client interface {
	SearchMetrics(ctx context.Context, actor domain.ActorContext, input contract.SearchMetricsInput) ([]domain.MetricCandidate, error)
	DescribeMetric(ctx context.Context, actor domain.ActorContext, input contract.DescribeMetricInput) (domain.MetricDescriptor, error)
	QueryRange(ctx context.Context, actor domain.ActorContext, input domain.QueryRangeRequest) (domain.QueryResult, error)
}
