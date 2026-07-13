// Package validation implements policies that must run after agent output and
// before a ChartSpec crosses the trusted backend boundary.
package validation

import (
	"context"
	"strings"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
)

type ChartSpecValidator struct{}

func (ChartSpecValidator) Validate(_ context.Context, request domain.AnalysisRequest, spec domain.ChartSpec) error {
	if strings.TrimSpace(spec.ID) == "" || strings.TrimSpace(spec.Title) == "" || strings.TrimSpace(spec.PromQL) == "" {
		return domain.NewError(domain.CodeQueryInvalid, "agent returned an incomplete ChartSpec")
	}
	if spec.DatasourceUID != request.Scope.DatasourceUID || spec.TimeRange != request.Scope.TimeRange {
		return domain.NewError(domain.CodeDatasourceForbidden, "agent changed the trusted analysis scope")
	}
	switch spec.Type {
	case domain.ChartTypeTimeseries, domain.ChartTypeStat, domain.ChartTypeTable:
		return nil
	default:
		return domain.NewError(domain.CodeQueryInvalid, "agent returned an unsupported chart type")
	}
}

var _ domain.ChartSpecValidator = ChartSpecValidator{}
