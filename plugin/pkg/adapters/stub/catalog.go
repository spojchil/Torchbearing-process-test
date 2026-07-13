// Package stub provides deterministic MS1 skeleton adapters. They prove the
// dependency graph and are replaced by Issues #5 and #6.
package stub

import (
	"context"
	"strings"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
)

type MetricCatalog struct{}

func (MetricCatalog) Search(_ context.Context, _ domain.AnalysisScope, text string, limit int) ([]domain.MetricCandidate, error) {
	if strings.TrimSpace(text) == "" || limit <= 0 {
		return nil, domain.NewError(domain.CodeInvalidRequest, "metric search input is invalid")
	}
	return []domain.MetricCandidate{
		{
			Name:  "http_requests_total",
			Type:  "counter",
			Help:  "Stub counter used to prove the MS1 architecture path",
			Score: 1,
		},
	}, nil
}

func (MetricCatalog) Describe(_ context.Context, _ domain.AnalysisScope, metric string) (domain.MetricDescriptor, error) {
	if metric != "http_requests_total" {
		return domain.MetricDescriptor{}, domain.NewError(domain.CodeMetricNotFound, "metric was not found")
	}
	return domain.MetricDescriptor{
		Name:   metric,
		Type:   "counter",
		Help:   "Stub counter used to prove the MS1 architecture path",
		Labels: []string{"service"},
	}, nil
}

var _ domain.MetricCatalog = MetricCatalog{}
