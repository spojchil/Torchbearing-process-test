package deterministic

import (
	"testing"
	"time"
)

// TestIDGeneratorIsDeterministic 验证 ID 格式和递增顺序固定。
func TestIDGeneratorIsDeterministic(t *testing.T) {
	generator := NewIDGenerator("")

	if got := generator.Next(); got != "mock-analysis-001" {
		t.Fatalf("first ID = %q, want mock-analysis-001", got)
	}
	if got := generator.Next(); got != "mock-analysis-002" {
		t.Fatalf("second ID = %q, want mock-analysis-002", got)
	}
}

// TestClockNeverDrifts 验证多次读取始终得到同一个注入时间。
func TestClockNeverDrifts(t *testing.T) {
	want := time.Date(2026, time.July, 13, 0, 0, 0, 0, time.UTC)
	clock := NewClock(want)

	if first, second := clock.Now(), clock.Now(); !first.Equal(want) || !second.Equal(want) {
		t.Fatalf("clock values = %v and %v, want %v", first, second, want)
	}
}

// TestScenarioValidation 验证未声明场景不会被 mock 静默接受。
func TestScenarioValidation(t *testing.T) {
	if !ScenarioBoundary.Valid() {
		t.Fatal("boundary scenario should be valid")
	}
	if Scenario("unknown").Valid() {
		t.Fatal("unknown scenario should be invalid")
	}
}
