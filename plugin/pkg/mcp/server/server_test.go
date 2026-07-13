package server

import (
	"context"
	"testing"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
	"github.com/spojchil/torchbearing/pkg/mcp/contract"
)

type recordingCatalog struct {
	searchCalls int
}

func (c *recordingCatalog) Search(context.Context, domain.AnalysisScope, string, int) ([]domain.MetricCandidate, error) {
	c.searchCalls++
	return []domain.MetricCandidate{{Name: "metric"}}, nil
}

func (*recordingCatalog) Describe(context.Context, domain.AnalysisScope, string) (domain.MetricDescriptor, error) {
	return domain.MetricDescriptor{Name: "metric"}, nil
}

type noopGateway struct{}

func (noopGateway) QueryRange(context.Context, domain.ActorContext, domain.QueryRangeRequest) (domain.QueryResult, error) {
	return domain.QueryResult{Status: domain.QueryStatusSuccess}, nil
}

func TestServerRejectsDatasourceBeforeDomainPort(t *testing.T) {
	catalog := &recordingCatalog{}
	server := New(catalog, noopGateway{})
	actor := domain.ActorContext{Access: domain.AccessScope{AllowedDatasourceUIDs: []string{"allowed"}}}

	_, err := server.SearchMetrics(context.Background(), actor, contract.SearchMetricsInput{
		Scope: domain.AnalysisScope{DatasourceUID: "forbidden"},
		Text:  "request rate",
		Limit: 5,
	})
	if domain.ErrorCodeOf(err) != domain.CodeDatasourceForbidden {
		t.Fatalf("expected DATASOURCE_FORBIDDEN, got %v", err)
	}
	if catalog.searchCalls != 0 {
		t.Fatalf("domain catalog was called %d times", catalog.searchCalls)
	}
}

func TestServerCallsDomainPortForAllowedDatasource(t *testing.T) {
	catalog := &recordingCatalog{}
	server := New(catalog, noopGateway{})
	actor := domain.ActorContext{Access: domain.AccessScope{AllowedDatasourceUIDs: []string{"allowed"}}}

	_, err := server.SearchMetrics(context.Background(), actor, contract.SearchMetricsInput{
		Scope: domain.AnalysisScope{DatasourceUID: "allowed"},
		Text:  "request rate",
		Limit: 5,
	})
	if err != nil {
		t.Fatalf("search metrics: %v", err)
	}
	if catalog.searchCalls != 1 {
		t.Fatalf("domain catalog was called %d times", catalog.searchCalls)
	}
}
