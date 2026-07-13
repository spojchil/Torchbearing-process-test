// Package metrics 定义 MS1 的指标 SDK 边界，具体指标适配器由其他模块负责实现。
package metrics

import (
	"context"

	"github.com/spojchil/torchbearing/internal/core"
)

// Client 执行强类型指标查询并返回强类型结果。
type Client interface {
	Query(ctx context.Context, query core.MetricQuery) (core.MetricResult, error)
}
