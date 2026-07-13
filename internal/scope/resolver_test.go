package scope

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spojchil/torchbearing/internal/contracts"
	"github.com/spojchil/torchbearing/internal/core"
)

// TestResolverScenarios 覆盖成功、空输入、可预期失败和相等时间边界。
func TestResolverScenarios(t *testing.T) {
	tests := []struct {
		name        string
		input       core.AnalysisScope
		want        core.AnalysisScope
		wantCode    core.ErrorCode
		wantMessage string
	}{
		{
			name: "success normalizes explicit scope",
			input: core.AnalysisScope{
				DatasourceUID: "  prometheus-mock  ",
				TimeRange:     core.TimeRange{From: " now-30m ", To: " now "},
			},
			want: core.AnalysisScope{
				DatasourceUID: "prometheus-mock",
				TimeRange:     core.TimeRange{From: "now-30m", To: "now"},
			},
		},
		{
			name:        "empty scope returns no context",
			input:       core.AnalysisScope{},
			wantCode:    core.ErrorCodeInvalidScope,
			wantMessage: "datasource UID is required",
		},
		{
			name: "predictable malformed-relative failure",
			input: core.AnalysisScope{
				DatasourceUID: "prometheus-mock",
				TimeRange:     core.TimeRange{From: "now-bad", To: "now"},
			},
			wantCode:    core.ErrorCodeInvalidScope,
			wantMessage: "invalid time range start",
		},
		{
			name: "equal absolute timestamps are an accepted boundary",
			input: core.AnalysisScope{
				DatasourceUID: "prometheus-mock",
				TimeRange: core.TimeRange{
					From: "2026-07-13T00:00:00Z",
					To:   "2026-07-13T00:00:00Z",
				},
			},
			want: core.AnalysisScope{
				DatasourceUID: "prometheus-mock",
				TimeRange: core.TimeRange{
					From: "2026-07-13T00:00:00Z",
					To:   "2026-07-13T00:00:00Z",
				},
			},
		},
	}

	resolver := NewResolver()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := resolver.Resolve(context.Background(), test.input)
			if test.wantCode == "" {
				if err != nil {
					t.Fatalf("Resolve returned unexpected error: %v", err)
				}
				if got.Scope != test.want {
					t.Fatalf("Scope = %+v, want %+v", got.Scope, test.want)
				}
				return
			}

			if got != (core.AnalysisContext{}) {
				t.Fatalf("context = %+v, want zero value", got)
			}
			assertAppError(t, err, test.wantCode, test.wantMessage)
		})
	}
}

// TestResolverMatchesInvalidScopeContractFixture 确保 B 的实现遵守 A 冻结的边界夹具。
func TestResolverMatchesInvalidScopeContractFixture(t *testing.T) {
	fixture := readInvalidScopeFixture(t)
	resolver := NewResolver()

	got, err := resolver.Resolve(context.Background(), core.AnalysisScope{
		DatasourceUID: fixture.Request.Scope.DatasourceUID,
		TimeRange: core.TimeRange{
			From: fixture.Request.Scope.TimeRange.From,
			To:   fixture.Request.Scope.TimeRange.To,
		},
	})

	if got != (core.AnalysisContext{}) {
		t.Fatalf("context = %+v, want zero value", got)
	}
	assertAppError(t, err, fixture.Error.Code, fixture.Error.Message)
}

// TestResolverReturnsTypedCancellation 验证取消原因仍可沿标准错误链识别。
func TestResolverReturnsTypedCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := NewResolver().Resolve(ctx, core.AnalysisScope{
		DatasourceUID: "prometheus-mock",
		TimeRange:     core.TimeRange{From: "now-30m", To: "now"},
	})

	if got != (core.AnalysisContext{}) {
		t.Fatalf("context = %+v, want zero value", got)
	}
	assertAppError(t, err, core.ErrorCodeInvalidArgument, "scope resolution canceled")
	if !errors.Is(err, context.Canceled) {
		t.Fatal("cancellation cause is not discoverable with errors.Is")
	}
}

// TestResolverPreservesOpaqueGrafanaExpressions 验证未解析表达式不会被 B 模块擅自修改。
func TestResolverPreservesOpaqueGrafanaExpressions(t *testing.T) {
	input := core.AnalysisScope{
		DatasourceUID: "prometheus-mock",
		TimeRange:     core.TimeRange{From: "now/d", To: "now/d+1d"},
	}

	got, err := NewResolver().Resolve(context.Background(), input)
	if err != nil {
		t.Fatalf("Resolve returned unexpected error: %v", err)
	}
	if got.Scope != input {
		t.Fatalf("Scope = %+v, want opaque scope unchanged: %+v", got.Scope, input)
	}
}

// assertAppError 统一校验 B 模块返回的强类型错误字段。
func assertAppError(t *testing.T, err error, code core.ErrorCode, message string) {
	t.Helper()

	var typed *core.AppError
	if !errors.As(err, &typed) {
		t.Fatalf("error type = %T, want *core.AppError", err)
	}
	if typed.Code != code {
		t.Fatalf("Code = %q, want %q", typed.Code, code)
	}
	if typed.Message != message {
		t.Fatalf("Message = %q, want %q", typed.Message, message)
	}
	if typed.Retryable {
		t.Fatal("Retryable = true, want false")
	}
}

// invalidScopeFixture 只声明本测试需要读取的 A 契约夹具字段。
type invalidScopeFixture struct {
	Request contracts.AnalysisRequest `json:"request"`
	Error   contracts.ErrorResponse   `json:"error"`
}

// readInvalidScopeFixture 从仓库中的确定性夹具读取反向时间场景。
func readInvalidScopeFixture(t *testing.T) invalidScopeFixture {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate resolver test source")
	}
	path := filepath.Join(filepath.Dir(file), "..", "..", "contracts", "fixtures", "invalid-scope.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	var fixture invalidScopeFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return fixture
}
