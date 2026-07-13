// Package contract defines the transport-neutral names and inputs for the
// three MS1 observability tools. Neither Client nor Server owns this contract.
package contract

import domain "github.com/spojchil/torchbearing/pkg/domain/analysis"

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
