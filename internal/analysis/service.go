// Package analysis 实现 MS1 的最小分析编排，不包含任何真实外部系统访问。
package analysis

import (
	"context"
	"errors"
	"strings"

	"github.com/spojchil/torchbearing/internal/core"
	agentsdk "github.com/spojchil/torchbearing/sdk/agent"
	metricssdk "github.com/spojchil/torchbearing/sdk/metrics"
)

// Service 只依赖 A 定义的端口和 SDK，按固定顺序编排 MS1 主流程。
type Service struct {
	scopes  core.ScopeResolver
	agent   agentsdk.Client
	metrics metricssdk.Client
	charts  core.ChartBuilder
	ids     core.IDGenerator
}

// NewService 注入主流程需要的全部抽象依赖。
func NewService(
	scopes core.ScopeResolver,
	agent agentsdk.Client,
	metrics metricssdk.Client,
	charts core.ChartBuilder,
	ids core.IDGenerator,
) *Service {
	return &Service{
		scopes:  scopes,
		agent:   agent,
		metrics: metrics,
		charts:  charts,
		ids:     ids,
	}
}

var _ core.Analyzer = (*Service)(nil)

// Analyze 依次完成范围解析、计划生成、指标查询和图表构建。
func (s *Service) Analyze(ctx context.Context, request core.AnalysisRequest) (core.AnalysisResponse, error) {
	requestID := s.ids.Next()
	text := strings.TrimSpace(request.Text)
	if text == "" {
		return core.AnalysisResponse{}, core.NewAppError(
			core.ErrorCodeInvalidArgument,
			"analysis text is required",
			false,
		).WithRequestID(requestID)
	}

	analysisContext, err := s.scopes.Resolve(ctx, request.Scope)
	if err != nil {
		return core.AnalysisResponse{}, withRequestID(err, requestID)
	}

	plan, err := s.agent.Plan(ctx, core.AgentRequest{Text: text, Context: analysisContext})
	if err != nil {
		return core.AnalysisResponse{}, withRequestID(err, requestID)
	}

	metricResult, err := s.metrics.Query(ctx, plan.Query)
	if err != nil {
		return core.AnalysisResponse{}, withRequestID(err, requestID)
	}

	chartSpecs, err := s.charts.Build(plan, metricResult)
	if err != nil {
		return core.AnalysisResponse{}, withRequestID(err, requestID)
	}
	if chartSpecs == nil {
		chartSpecs = []core.ChartSpec{}
	}

	return core.AnalysisResponse{
		RequestID: requestID,
		Message:   plan.Message,
		Charts:    chartSpecs,
		Mock:      true,
	}, nil
}

func withRequestID(err error, requestID string) error {
	var typed *core.AppError
	if errors.As(err, &typed) {
		return typed.WithRequestID(requestID)
	}
	return core.WrapAppError(
		core.ErrorCodeInternal,
		"analysis failed unexpectedly",
		false,
		err,
	).WithRequestID(requestID)
}
