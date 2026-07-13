package core

// ChartType 表示分析结果返回的、与具体渲染器无关的图表类型。
type ChartType string

const (
	ChartTypeTimeSeries ChartType = "timeseries"
	ChartTypeStat       ChartType = "stat"
	ChartTypeTable      ChartType = "table"
)

// DataState 表示指标查询是否返回了可用数据。
type DataState string

const (
	DataStatePresent DataState = "present"
	DataStateEmpty   DataState = "empty"
)

// TimeRange 表示兼容 Grafana 的不透明时间范围，解析和校验由 scope 模块负责。
type TimeRange struct {
	From string
	To   string
}

// AnalysisScope 表示调用方提供的完整 MS1 分析范围。
type AnalysisScope struct {
	DatasourceUID string
	TimeRange     TimeRange
}

// AnalysisRequest 表示分析用例的领域输入。
type AnalysisRequest struct {
	Text  string
	Scope AnalysisScope
}

// AnalysisContext 表示经过校验、可供下游 SDK 使用的分析上下文。
type AnalysisContext struct {
	Scope AnalysisScope
}

// AgentRequest 表示 Agent SDK 的强类型输入。
type AgentRequest struct {
	Text    string
	Context AnalysisContext
}

// MetricQuery 表示 Metrics SDK 的强类型查询输入。
type MetricQuery struct {
	DatasourceUID string
	PromQL        string
	TimeRange     TimeRange
}

// ChartHint 表示 Agent 选择的展示意图，但不绑定具体前端图表库。
type ChartHint struct {
	ID    string
	Title string
	Type  ChartType
}

// AnalysisPlan 表示 Agent SDK 的强类型输出。
type AnalysisPlan struct {
	Message string
	Query   MetricQuery
	Chart   ChartHint
}

// Label 表示确定性的指标键值标签。这里使用切片而非 map，确保测试和序列化输出顺序稳定。
type Label struct {
	Name  string
	Value string
}

// DataPoint 表示一个带时间戳的指标值。
type DataPoint struct {
	Timestamp string
	Value     float64
}

// Series 表示一组具名指标序列。
type Series struct {
	Name   string
	Labels []Label
	Points []DataPoint
}

// MetricResult 表示 Metrics SDK 的强类型输出。
type MetricResult struct {
	State  DataState
	Series []Series
}

// ChartSpec 表示稳定且与渲染器无关的图表契约。
type ChartSpec struct {
	ID            string
	Title         string
	Type          ChartType
	DatasourceUID string
	PromQL        string
	TimeRange     TimeRange
}

// AnalysisResponse 表示分析用例的领域输出。
type AnalysisResponse struct {
	RequestID string
	Message   string
	Charts    []ChartSpec
	Mock      bool
}
