package contract_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spojchil/torchbearing/internal/contracts"
	"github.com/spojchil/torchbearing/internal/core"
)

// fixture 表示契约测试使用的成功或失败场景；响应与错误必须二选一。
type fixture struct {
	Scenario string                      `json:"scenario"`
	Request  contracts.AnalysisRequest   `json:"request"`
	Response *contracts.AnalysisResponse `json:"response,omitempty"`
	Error    *contracts.ErrorResponse    `json:"error,omitempty"`
}

// TestContractSchemasAreValidJSON 验证两份公共 Schema 至少具有合法 JSON 语法。
func TestContractSchemasAreValidJSON(t *testing.T) {
	for _, name := range []string{"analysis.schema.json", "error.schema.json"} {
		path := filepath.Join(repositoryRoot(t), "contracts", name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if !json.Valid(data) {
			t.Fatalf("%s is not valid JSON", path)
		}
	}
}

// TestDeterministicContractFixtures 验证所有 mock 场景的固定请求 ID 和结果类型。
func TestDeterministicContractFixtures(t *testing.T) {
	tests := []struct {
		file           string
		scenario       string
		requestID      string
		responseCharts int
		errorCode      core.ErrorCode
		retryable      bool
	}{
		{file: "success.json", scenario: "success", requestID: "mock-analysis-001", responseCharts: 1},
		{file: "empty.json", scenario: "empty", requestID: "mock-analysis-002", responseCharts: 0},
		{file: "agent-failure.json", scenario: "agent-failure", requestID: "mock-analysis-003", errorCode: core.ErrorCodeAgentUnavailable, retryable: true},
		{file: "metrics-failure.json", scenario: "metrics-failure", requestID: "mock-analysis-004", errorCode: core.ErrorCodeMetricsUnavailable, retryable: true},
		{file: "invalid-scope.json", scenario: "invalid-scope", requestID: "mock-analysis-005", errorCode: core.ErrorCodeInvalidScope, retryable: false},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			got := readFixture(t, test.file)
			if got.Scenario != test.scenario {
				t.Fatalf("Scenario = %q, want %q", got.Scenario, test.scenario)
			}
			if got.Request.Text == "" || got.Request.Scope.DatasourceUID == "" {
				t.Fatal("request is missing typed text or datasource scope")
			}

			if test.errorCode == "" {
				if got.Response == nil || got.Error != nil {
					t.Fatal("success fixture must contain response and no error")
				}
				if got.Response.RequestID != test.requestID {
					t.Fatalf("RequestID = %q, want %q", got.Response.RequestID, test.requestID)
				}
				if len(got.Response.Charts) != test.responseCharts {
					t.Fatalf("len(Charts) = %d, want %d", len(got.Response.Charts), test.responseCharts)
				}
				if !got.Response.Mock {
					t.Fatal("Mock = false, want true")
				}
				return
			}

			if got.Response != nil || got.Error == nil {
				t.Fatal("failure fixture must contain error and no response")
			}
			if got.Error.Code != test.errorCode {
				t.Fatalf("Code = %q, want %q", got.Error.Code, test.errorCode)
			}
			if got.Error.RequestID != test.requestID {
				t.Fatalf("RequestID = %q, want %q", got.Error.RequestID, test.requestID)
			}
			if got.Error.Retryable != test.retryable {
				t.Fatalf("Retryable = %v, want %v", got.Error.Retryable, test.retryable)
			}
		})
	}
}

// readFixture 读取并解码一份强类型契约夹具。
func readFixture(t *testing.T, name string) fixture {
	t.Helper()

	path := filepath.Join(repositoryRoot(t), "contracts", "fixtures", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	var got fixture
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return got
}

// repositoryRoot 根据当前测试文件位置确定仓库根目录，避免依赖执行目录。
func repositoryRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate contract test source")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
