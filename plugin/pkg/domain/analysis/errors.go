package analysis

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	CodeInvalidRequest      ErrorCode = "INVALID_REQUEST"
	CodeIntentAmbiguous     ErrorCode = "INTENT_AMBIGUOUS"
	CodeMetricNotFound      ErrorCode = "METRIC_NOT_FOUND"
	CodeQueryInvalid        ErrorCode = "QUERY_INVALID"
	CodeNoData              ErrorCode = "NO_DATA"
	CodeQueryLimited        ErrorCode = "QUERY_LIMITED"
	CodeDatasourceForbidden ErrorCode = "DATASOURCE_FORBIDDEN"
	CodeModelUnavailable    ErrorCode = "MODEL_UNAVAILABLE"
	CodeMCPUnavailable      ErrorCode = "MCP_UNAVAILABLE"
	CodeInternal            ErrorCode = "INTERNAL"
)

type CodedError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *CodedError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

func (e *CodedError) Unwrap() error {
	return e.Cause
}

func NewError(code ErrorCode, message string) error {
	return &CodedError{Code: code, Message: message}
}

func WrapError(code ErrorCode, message string, cause error) error {
	return &CodedError{Code: code, Message: message, Cause: cause}
}

func ErrorCodeOf(err error) ErrorCode {
	var coded *CodedError
	if errors.As(err, &coded) {
		return coded.Code
	}
	return CodeInternal
}

func SafeMessage(err error) string {
	var coded *CodedError
	if errors.As(err, &coded) && coded.Message != "" {
		return coded.Message
	}
	return "internal error"
}

func RequireDatasourceAccess(actor ActorContext, datasourceUID string) error {
	if actor.Access.AllowsDatasource(datasourceUID) {
		return nil
	}
	return NewError(CodeDatasourceForbidden, fmt.Sprintf("datasource %q is not allowed", datasourceUID))
}
