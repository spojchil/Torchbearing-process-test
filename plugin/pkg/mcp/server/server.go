// Package server contains the transport-independent handlers behind the three
// MS1 observability MCP tools. A Streamable HTTP adapter can wrap this package
// without moving domain behavior into the protocol layer.
package server

import (
	"context"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
	"github.com/spojchil/torchbearing/pkg/mcp/contract"
)

type Server struct {
	catalog domain.MetricCatalog
	gateway domain.PrometheusGateway
}

func New(catalog domain.MetricCatalog, gateway domain.PrometheusGateway) *Server {
	return &Server{catalog: catalog, gateway: gateway}
}

func (s *Server) SearchMetrics(ctx context.Context, actor domain.ActorContext, input contract.SearchMetricsInput) ([]domain.MetricCandidate, error) {
	if err := domain.RequireDatasourceAccess(actor, input.Scope.DatasourceUID); err != nil {
		return nil, err
	}
	return s.catalog.Search(ctx, input.Scope, input.Text, input.Limit)
}

func (s *Server) DescribeMetric(ctx context.Context, actor domain.ActorContext, input contract.DescribeMetricInput) (domain.MetricDescriptor, error) {
	if err := domain.RequireDatasourceAccess(actor, input.Scope.DatasourceUID); err != nil {
		return domain.MetricDescriptor{}, err
	}
	return s.catalog.Describe(ctx, input.Scope, input.Metric)
}

func (s *Server) QueryRange(ctx context.Context, actor domain.ActorContext, input domain.QueryRangeRequest) (domain.QueryResult, error) {
	if err := domain.RequireDatasourceAccess(actor, input.DatasourceUID); err != nil {
		return domain.QueryResult{}, err
	}
	return s.gateway.QueryRange(ctx, actor, input)
}
