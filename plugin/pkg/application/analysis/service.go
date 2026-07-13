// Package analysis implements the use-case orchestration between the trusted
// plugin boundary, the agent port, and domain policy validation.
package analysis

import (
	"context"
	"strings"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
)

type AgentRunner interface {
	Analyze(ctx context.Context, actor domain.ActorContext, request domain.AnalysisRequest) (domain.AnalysisResult, error)
}

type Service struct {
	runner    AgentRunner
	validator domain.ChartSpecValidator
}

func NewService(runner AgentRunner, validator domain.ChartSpecValidator) *Service {
	return &Service{runner: runner, validator: validator}
}

func (s *Service) Analyze(ctx context.Context, actor domain.ActorContext, request domain.AnalysisRequest) (domain.AnalysisResult, error) {
	if strings.TrimSpace(request.Text) == "" || strings.TrimSpace(request.Scope.DatasourceUID) == "" {
		return domain.AnalysisResult{}, domain.NewError(domain.CodeInvalidRequest, "text and datasourceUid are required")
	}
	if request.Scope.TimeRange.From == "" || request.Scope.TimeRange.To == "" {
		return domain.AnalysisResult{}, domain.NewError(domain.CodeInvalidRequest, "time range is required")
	}
	if err := domain.RequireDatasourceAccess(actor, request.Scope.DatasourceUID); err != nil {
		return domain.AnalysisResult{}, err
	}

	result, err := s.runner.Analyze(ctx, actor, request)
	if err != nil {
		return domain.AnalysisResult{}, err
	}
	if err := s.validator.Validate(ctx, request, result.Chart); err != nil {
		return domain.AnalysisResult{}, err
	}

	return result, nil
}
