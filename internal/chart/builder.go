// Package chart 将强类型分析计划转换为与前端渲染器无关的 MS1 图表定义。
package chart

import (
	"strings"

	"github.com/spojchil/torchbearing/internal/core"
)

// Builder 是不依赖 Grafana 或第三方图表库的确定性图表构建器。
type Builder struct{}

// NewBuilder 创建无外部依赖的图表构建器。
func NewBuilder() *Builder {
	return &Builder{}
}

var _ core.ChartBuilder = (*Builder)(nil)

// Build 根据指标数据状态返回一个图表、空图表集合或强类型错误。
func (b *Builder) Build(plan core.AnalysisPlan, result core.MetricResult) ([]core.ChartSpec, error) {
	switch result.State {
	case core.DataStateEmpty:
		return []core.ChartSpec{}, nil
	case core.DataStatePresent:
		if len(result.Series) == 0 {
			return nil, core.NewAppError(
				core.ErrorCodeNoData,
				"present metrics result must contain at least one series",
				false,
			)
		}
	default:
		return nil, core.NewAppError(
			core.ErrorCodeInvalidArgument,
			"metrics data state is invalid",
			false,
		)
	}

	if err := validatePlan(plan); err != nil {
		return nil, err
	}

	return []core.ChartSpec{
		{
			ID:            plan.Chart.ID,
			Title:         plan.Chart.Title,
			Type:          plan.Chart.Type,
			DatasourceUID: plan.Query.DatasourceUID,
			PromQL:        plan.Query.PromQL,
			TimeRange:     plan.Query.TimeRange,
		},
	}, nil
}

func validatePlan(plan core.AnalysisPlan) error {
	if strings.TrimSpace(plan.Chart.ID) == "" || strings.TrimSpace(plan.Chart.Title) == "" {
		return core.NewAppError(core.ErrorCodeInvalidArgument, "chart identity is required", false)
	}
	switch plan.Chart.Type {
	case core.ChartTypeTimeSeries, core.ChartTypeStat, core.ChartTypeTable:
	default:
		return core.NewAppError(core.ErrorCodeInvalidArgument, "chart type is invalid", false)
	}
	if strings.TrimSpace(plan.Query.DatasourceUID) == "" || strings.TrimSpace(plan.Query.PromQL) == "" {
		return core.NewAppError(core.ErrorCodeInvalidArgument, "chart query is required", false)
	}
	if strings.TrimSpace(plan.Query.TimeRange.From) == "" || strings.TrimSpace(plan.Query.TimeRange.To) == "" {
		return core.NewAppError(core.ErrorCodeInvalidArgument, "chart time range is required", false)
	}
	return nil
}
