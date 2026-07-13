package core

import (
	"context"
	"time"
)

// ScopeResolver 负责校验并规范化调用方提供的分析范围。
type ScopeResolver interface {
	Resolve(ctx context.Context, scope AnalysisScope) (AnalysisContext, error)
}

// ChartBuilder 将 Agent 计划和指标结果转换为与渲染器无关的图表定义。
type ChartBuilder interface {
	Build(plan AnalysisPlan, result MetricResult) ([]ChartSpec, error)
}

// IDGenerator 为编排层提供确定性的请求标识。
type IDGenerator interface {
	Next() string
}

// Clock 隔离编排层对时间的依赖，使 MS1 适配器和测试保持确定性。
type Clock interface {
	Now() time.Time
}

// Analyzer 是 MS1 暴露给上层的唯一分析用例边界。
type Analyzer interface {
	Analyze(ctx context.Context, request AnalysisRequest) (AnalysisResponse, error)
}
