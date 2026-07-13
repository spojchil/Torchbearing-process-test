// Package analysis contains the MS1 runtime model and domain ports. It has no
// dependency on Grafana, HTTP, MCP, or an agent framework.
package analysis

import "time"

type AccessScope struct {
	AllowedDatasourceUIDs []string
}

func (s AccessScope) AllowsDatasource(uid string) bool {
	for _, allowed := range s.AllowedDatasourceUIDs {
		if uid == allowed {
			return true
		}
	}
	return false
}

type ActorContext struct {
	OrgID  int64
	Login  string
	Role   string
	Access AccessScope
}

type TimeRange struct {
	From string
	To   string
}

type AnalysisScope struct {
	DatasourceUID string
	TimeRange     TimeRange
	Service       string
	Environment   string
	Namespace     string
	Cluster       string
}

type AnalysisRequest struct {
	Text  string
	Scope AnalysisScope
}

type MetricCandidate struct {
	Name  string
	Type  string
	Help  string
	Score float64
}

type MetricDescriptor struct {
	Name   string
	Type   string
	Help   string
	Labels []string
}

type QueryRangeRequest struct {
	DatasourceUID string
	Expression    string
	TimeRange     TimeRange
	Step          time.Duration
	MaxSeries     int
}

type QueryStatus string

const (
	QueryStatusSuccess QueryStatus = "success"
	QueryStatusNoData  QueryStatus = "no_data"
)

type QueryResult struct {
	Status      QueryStatus
	Duration    time.Duration
	SeriesCount int
}

type ChartType string

const (
	ChartTypeTimeseries ChartType = "timeseries"
	ChartTypeStat       ChartType = "stat"
	ChartTypeTable      ChartType = "table"
)

type ChartSpec struct {
	ID            string
	Title         string
	Type          ChartType
	DatasourceUID string
	PromQL        string
	TimeRange     TimeRange
	Unit          string
	Legend        string
}

type QuerySummary struct {
	Language    string
	Expression  string
	Status      QueryStatus
	Duration    time.Duration
	SeriesCount int
}

type Evidence struct {
	Metrics     []string
	Explanation string
}

type AnalysisResult struct {
	Chart    ChartSpec
	Query    QuerySummary
	Evidence Evidence
}
