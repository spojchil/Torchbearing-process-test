package core

import (
	"errors"
	"testing"
)

// TestAppErrorPreservesTypedFieldsAndCause 验证强类型字段及底层错误链不会丢失。
func TestAppErrorPreservesTypedFieldsAndCause(t *testing.T) {
	cause := errors.New("deterministic agent failure")
	err := WrapAppError(
		ErrorCodeAgentUnavailable,
		"agent is unavailable",
		true,
		cause,
	).WithRequestID("mock-analysis-003")

	if err.Code != ErrorCodeAgentUnavailable {
		t.Fatalf("Code = %q, want %q", err.Code, ErrorCodeAgentUnavailable)
	}
	if !err.Retryable {
		t.Fatal("Retryable = false, want true")
	}
	if err.RequestID != "mock-analysis-003" {
		t.Fatalf("RequestID = %q, want mock-analysis-003", err.RequestID)
	}
	if !errors.Is(err, cause) {
		t.Fatal("wrapped cause is not discoverable with errors.Is")
	}

	var typed *AppError
	if !errors.As(err, &typed) {
		t.Fatal("error is not discoverable as *AppError")
	}
}

// TestAppErrorWithRequestIDDoesNotMutateOriginal 验证追加请求 ID 时采用复制语义。
func TestAppErrorWithRequestIDDoesNotMutateOriginal(t *testing.T) {
	original := NewAppError(ErrorCodeInvalidScope, "invalid time range", false)
	annotated := original.WithRequestID("mock-analysis-005")

	if original.RequestID != "" {
		t.Fatalf("original RequestID = %q, want empty", original.RequestID)
	}
	if annotated.RequestID != "mock-analysis-005" {
		t.Fatalf("annotated RequestID = %q, want mock-analysis-005", annotated.RequestID)
	}
}
