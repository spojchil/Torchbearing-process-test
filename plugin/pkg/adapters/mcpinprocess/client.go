// Package mcpinprocess is the MS1 skeleton transport. Issue #5 replaces it
// with Streamable HTTP while preserving the mcp/client port.
package mcpinprocess

import (
	"context"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
	mcpclient "github.com/spojchil/torchbearing/pkg/mcp/client"
	"github.com/spojchil/torchbearing/pkg/mcp/contract"
	mcpserver "github.com/spojchil/torchbearing/pkg/mcp/server"
)

type Client struct {
	server *mcpserver.Server
}

func New(server *mcpserver.Server) *Client {
	return &Client{server: server}
}

func (c *Client) SearchMetrics(ctx context.Context, actor domain.ActorContext, input contract.SearchMetricsInput) ([]domain.MetricCandidate, error) {
	return c.server.SearchMetrics(ctx, actor, input)
}

func (c *Client) DescribeMetric(ctx context.Context, actor domain.ActorContext, input contract.DescribeMetricInput) (domain.MetricDescriptor, error) {
	return c.server.DescribeMetric(ctx, actor, input)
}

func (c *Client) QueryRange(ctx context.Context, actor domain.ActorContext, input domain.QueryRangeRequest) (domain.QueryResult, error) {
	return c.server.QueryRange(ctx, actor, input)
}

var _ mcpclient.Client = (*Client)(nil)
