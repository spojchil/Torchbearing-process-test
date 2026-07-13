package stub

import (
	"context"
	"strings"
	"time"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
)

type PrometheusGateway struct{}

func (PrometheusGateway) QueryRange(_ context.Context, actor domain.ActorContext, request domain.QueryRangeRequest) (domain.QueryResult, error) {
	if err := domain.RequireDatasourceAccess(actor, request.DatasourceUID); err != nil {
		return domain.QueryResult{}, err
	}
	if strings.TrimSpace(request.Expression) == "" {
		return domain.QueryResult{}, domain.NewError(domain.CodeQueryInvalid, "PromQL expression is required")
	}
	return domain.QueryResult{
		Status:      domain.QueryStatusSuccess,
		Duration:    12 * time.Millisecond,
		SeriesCount: 1,
	}, nil
}

var _ domain.PrometheusGateway = PrometheusGateway{}
