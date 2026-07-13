package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spojchil/torchbearing/internal/contracts"
)

// TestRunOutputsDeterministicSuccess 验证命令入口输出稳定成功响应。
func TestRunOutputsDeterministicSuccess(t *testing.T) {
	var output bytes.Buffer
	if err := run(&output); err != nil {
		t.Fatalf("run returned unexpected error: %v", err)
	}

	var response contracts.AnalysisResponse
	if err := json.Unmarshal(output.Bytes(), &response); err != nil {
		t.Fatalf("decode command output: %v", err)
	}
	if response.RequestID != "mock-analysis-001" {
		t.Fatalf("request ID = %q, want mock-analysis-001", response.RequestID)
	}
	if !response.Mock || len(response.Charts) != 1 {
		t.Fatalf("response = %+v, want one mock chart", response)
	}
}
