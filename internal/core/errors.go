package core

import "fmt"

// ErrorCode 表示 MS1 各模块共享的稳定错误码。
type ErrorCode string

const (
	ErrorCodeInvalidArgument    ErrorCode = "INVALID_ARGUMENT"
	ErrorCodeInvalidScope       ErrorCode = "INVALID_SCOPE"
	ErrorCodeAgentUnavailable   ErrorCode = "AGENT_UNAVAILABLE"
	ErrorCodeMetricsUnavailable ErrorCode = "METRICS_UNAVAILABLE"
	ErrorCodeNoData             ErrorCode = "NO_DATA"
	ErrorCodeInternal           ErrorCode = "INTERNAL"
)

// AppError 表示 MS1 的强类型错误。底层原因保持私有，避免传输层意外暴露适配器细节。
type AppError struct {
	Code      ErrorCode
	Message   string
	Retryable bool
	RequestID string
	cause     error
}

// NewAppError 创建不包含底层原因的强类型错误。
func NewAppError(code ErrorCode, message string, retryable bool) *AppError {
	return &AppError{Code: code, Message: message, Retryable: retryable}
}

// WrapAppError 创建强类型错误，同时保留底层原因以支持 errors.Is 和 errors.As。
func WrapAppError(code ErrorCode, message string, retryable bool, cause error) *AppError {
	return &AppError{Code: code, Message: message, Retryable: retryable, cause: cause}
}

// WithRequestID 返回带公开请求标识的错误副本，不修改原错误。
func (e *AppError) WithRequestID(requestID string) *AppError {
	if e == nil {
		return nil
	}

	copy := *e
	copy.RequestID = requestID
	return &copy
}

func (e *AppError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 只暴露标准错误链，不暴露传输数据。
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.cause
}
