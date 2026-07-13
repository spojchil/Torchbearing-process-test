package analysis

import "context"

type MetricCatalog interface {
	Search(ctx context.Context, scope AnalysisScope, text string, limit int) ([]MetricCandidate, error)
	Describe(ctx context.Context, scope AnalysisScope, metric string) (MetricDescriptor, error)
}

type PrometheusGateway interface {
	QueryRange(ctx context.Context, actor ActorContext, request QueryRangeRequest) (QueryResult, error)
}

type ChartSpecValidator interface {
	Validate(ctx context.Context, request AnalysisRequest, spec ChartSpec) error
}
