// Package bootstrap 负责组装 MS1 的确定性依赖，不承载业务逻辑。
package bootstrap

import (
	"github.com/spojchil/torchbearing/internal/analysis"
	"github.com/spojchil/torchbearing/internal/chart"
	"github.com/spojchil/torchbearing/internal/core"
	"github.com/spojchil/torchbearing/internal/scope"
	"github.com/spojchil/torchbearing/internal/transport/grafana"
	agentmock "github.com/spojchil/torchbearing/mocks/agent"
	"github.com/spojchil/torchbearing/mocks/deterministic"
	metricsmock "github.com/spojchil/torchbearing/mocks/metrics"
)

// NewMS1Analyzer 使用 B/C 的确定性实现组装 A 定义的 Analyzer 接口。
func NewMS1Analyzer(scenario deterministic.Scenario) (core.Analyzer, error) {
	if !scenario.Valid() {
		return nil, core.NewAppError(
			core.ErrorCodeInvalidArgument,
			"unsupported MS1 bootstrap scenario",
			false,
		)
	}

	ids := deterministic.NewIDGenerator("mock-analysis")
	// 预推进序号，使每个独立场景与 A 冻结的 fixture request ID 保持一致。
	for index := 0; index < scenarioOffset(scenario); index++ {
		ids.Next()
	}

	return analysis.NewService(
		scope.NewResolver(),
		agentmock.NewClient(scenario),
		metricsmock.NewClient(scenario),
		chart.NewBuilder(),
		ids,
	), nil
}

// NewMS1Gateway 创建可供入口或测试调用的确定性内存网关。
func NewMS1Gateway(scenario deterministic.Scenario) (*grafana.Gateway, error) {
	analyzer, err := NewMS1Analyzer(scenario)
	if err != nil {
		return nil, err
	}
	return grafana.NewGateway(analyzer), nil
}

func scenarioOffset(scenario deterministic.Scenario) int {
	switch scenario {
	case deterministic.ScenarioSuccess:
		return 0
	case deterministic.ScenarioEmpty:
		return 1
	case deterministic.ScenarioAgentFailure:
		return 2
	case deterministic.ScenarioMetricsFailure:
		return 3
	case deterministic.ScenarioBoundary:
		return 4
	default:
		return 0
	}
}
