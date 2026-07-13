package contracts

import "github.com/spojchil/torchbearing/internal/core"

// ErrorResponse 表示稳定的 MS1 错误响应契约。
type ErrorResponse struct {
	Code      core.ErrorCode `json:"code"`
	Message   string         `json:"message"`
	Retryable bool           `json:"retryable"`
	RequestID string         `json:"requestId"`
}
