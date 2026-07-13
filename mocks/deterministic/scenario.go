// Package deterministic 提供 MS1 mock 共用的确定性场景和基础适配器。
package deterministic

// Scenario 表示 Agent 与 Metrics mock 共享的固定执行场景。
type Scenario string

const (
	ScenarioSuccess        Scenario = "success"
	ScenarioEmpty          Scenario = "empty"
	ScenarioAgentFailure   Scenario = "agent-failure"
	ScenarioMetricsFailure Scenario = "metrics-failure"
	ScenarioBoundary       Scenario = "boundary"
)

// Valid 判断场景是否属于 C 模块支持的 MS1 确定性集合。
func (s Scenario) Valid() bool {
	switch s {
	case ScenarioSuccess,
		ScenarioEmpty,
		ScenarioAgentFailure,
		ScenarioMetricsFailure,
		ScenarioBoundary:
		return true
	default:
		return false
	}
}
