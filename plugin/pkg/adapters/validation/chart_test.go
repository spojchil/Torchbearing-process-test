package validation

import (
	"context"
	"testing"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
)

func TestChartSpecValidatorRejectsScopeChange(t *testing.T) {
	request := domain.AnalysisRequest{
		Scope: domain.AnalysisScope{
			DatasourceUID: "allowed",
			TimeRange:     domain.TimeRange{From: "now-30m", To: "now"},
		},
	}
	spec := domain.ChartSpec{
		ID:            "chart",
		Title:         "Chart",
		Type:          domain.ChartTypeTimeseries,
		DatasourceUID: "other",
		PromQL:        "up",
		TimeRange:     request.Scope.TimeRange,
	}

	err := (ChartSpecValidator{}).Validate(context.Background(), request, spec)
	if domain.ErrorCodeOf(err) != domain.CodeDatasourceForbidden {
		t.Fatalf("expected DATASOURCE_FORBIDDEN, got %v", err)
	}
}
