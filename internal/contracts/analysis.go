package contracts

import "github.com/spojchil/torchbearing/internal/core"

// TimeRange 表示分析时间范围的传输结构。
type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// AnalysisScope 表示调用方所选分析范围的传输结构。
type AnalysisScope struct {
	DatasourceUID string    `json:"datasourceUid"`
	TimeRange     TimeRange `json:"timeRange"`
}

// AnalysisRequest 表示稳定的 MS1 请求契约。
type AnalysisRequest struct {
	Text  string        `json:"text"`
	Scope AnalysisScope `json:"scope"`
}

// ChartSpec 表示与具体渲染器无关的图表定义。
type ChartSpec struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Type          core.ChartType `json:"type"`
	DatasourceUID string         `json:"datasourceUid"`
	PromQL        string         `json:"promql"`
	TimeRange     TimeRange      `json:"timeRange"`
}

// AnalysisResponse 表示稳定的 MS1 成功响应契约。
type AnalysisResponse struct {
	RequestID string      `json:"requestId"`
	Message   string      `json:"message"`
	Charts    []ChartSpec `json:"charts"`
	Mock      bool        `json:"mock"`
}
