// Package scope 校验并规范化调用方显式提供的分析范围，不依赖任何外部系统。
package scope

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/spojchil/torchbearing/internal/core"
)

// Resolver 是 MS1 使用的确定性内存范围解析器。
type Resolver struct{}

// NewResolver 创建不依赖外部能力的确定性范围解析器。
func NewResolver() *Resolver {
	return &Resolver{}
}

var _ core.ScopeResolver = (*Resolver)(nil)

// Resolve 清理稳定字符串字段，并校验 MS1 无需访问真实服务或系统时钟即可比较的时间范围。
// 其他非空 Grafana 时间表达式按照契约保持不透明，不在本模块中擅自改写。
func (r *Resolver) Resolve(ctx context.Context, input core.AnalysisScope) (core.AnalysisContext, error) {
	if err := ctx.Err(); err != nil {
		return core.AnalysisContext{}, core.WrapAppError(
			core.ErrorCodeInvalidArgument,
			"scope resolution canceled",
			false,
			err,
		)
	}

	normalized := core.AnalysisScope{
		DatasourceUID: strings.TrimSpace(input.DatasourceUID),
		TimeRange: core.TimeRange{
			From: strings.TrimSpace(input.TimeRange.From),
			To:   strings.TrimSpace(input.TimeRange.To),
		},
	}

	if normalized.DatasourceUID == "" {
		return invalidScope("datasource UID is required")
	}
	if normalized.TimeRange.From == "" {
		return invalidScope("time range start is required")
	}
	if normalized.TimeRange.To == "" {
		return invalidScope("time range end is required")
	}
	if err := validateTimeOrder(normalized.TimeRange); err != nil {
		return core.AnalysisContext{}, err
	}

	return core.AnalysisContext{Scope: normalized}, nil
}

func invalidScope(message string) (core.AnalysisContext, error) {
	return core.AnalysisContext{}, core.NewAppError(core.ErrorCodeInvalidScope, message, false)
}

type comparableTime struct {
	kind     timeKind
	absolute time.Time
	offset   time.Duration
}

type timeKind uint8

const (
	timeKindOpaque timeKind = iota
	timeKindAbsolute
	timeKindRelative
)

func validateTimeOrder(timeRange core.TimeRange) error {
	// 只有两端属于同一种可比较时间时才判断先后顺序；混合或不透明表达式交由后续适配器解释。
	from, err := parseComparableTime(timeRange.From)
	if err != nil {
		return core.NewAppError(core.ErrorCodeInvalidScope, "invalid time range start", false)
	}
	to, err := parseComparableTime(timeRange.To)
	if err != nil {
		return core.NewAppError(core.ErrorCodeInvalidScope, "invalid time range end", false)
	}

	if from.kind != to.kind || from.kind == timeKindOpaque {
		return nil
	}

	var startsAfterEnd bool
	switch from.kind {
	case timeKindAbsolute:
		startsAfterEnd = from.absolute.After(to.absolute)
	case timeKindRelative:
		startsAfterEnd = from.offset > to.offset
	}
	if startsAfterEnd {
		return core.NewAppError(
			core.ErrorCodeInvalidScope,
			"time range start must not be after end",
			false,
		)
	}
	return nil
}

func parseComparableTime(value string) (comparableTime, error) {
	if absolute, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return comparableTime{kind: timeKindAbsolute, absolute: absolute}, nil
	}

	offset, relative, err := parseRelativeTime(value)
	if err != nil {
		return comparableTime{}, err
	}
	if relative {
		return comparableTime{kind: timeKindRelative, offset: offset}, nil
	}

	return comparableTime{kind: timeKindOpaque}, nil
}

func parseRelativeTime(value string) (time.Duration, bool, error) {
	// MS1 仅解析无需真实时钟即可比较的 now 相对表达式。
	if value == "now" {
		return 0, true, nil
	}
	if !strings.HasPrefix(value, "now-") && !strings.HasPrefix(value, "now+") {
		return 0, false, nil
	}
	if len(value) < len("now-1s") {
		return 0, true, errors.New("relative time is missing an amount or unit")
	}

	amount, err := strconv.ParseInt(value[4:len(value)-1], 10, 64)
	if err != nil || amount < 0 {
		return 0, true, errors.New("relative time amount is invalid")
	}

	unit, err := relativeUnit(value[len(value)-1])
	if err != nil {
		return 0, true, err
	}
	if amount > math.MaxInt64/int64(unit) {
		return 0, true, errors.New("relative time is too large")
	}

	offset := time.Duration(amount) * unit
	if value[3] == '-' {
		offset = -offset
	}
	return offset, true, nil
}

func relativeUnit(unit byte) (time.Duration, error) {
	// 单位集合覆盖 MS1 使用的秒、分、时、天和周，不提前实现更复杂的日历语义。
	switch unit {
	case 's':
		return time.Second, nil
	case 'm':
		return time.Minute, nil
	case 'h':
		return time.Hour, nil
	case 'd':
		return 24 * time.Hour, nil
	case 'w':
		return 7 * 24 * time.Hour, nil
	default:
		return 0, errors.New("relative time unit is invalid")
	}
}
