package analysis

import (
	"context"
	"errors"
	"testing"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
)

type runnerFunc func(context.Context, domain.ActorContext, domain.AnalysisRequest) (domain.AnalysisResult, error)

func (f runnerFunc) Analyze(ctx context.Context, actor domain.ActorContext, request domain.AnalysisRequest) (domain.AnalysisResult, error) {
	return f(ctx, actor, request)
}

type validatorFunc func(context.Context, domain.AnalysisRequest, domain.ChartSpec) error

func (f validatorFunc) Validate(ctx context.Context, request domain.AnalysisRequest, spec domain.ChartSpec) error {
	return f(ctx, request, spec)
}

func TestServiceRejectsDatasourceBeforeAgent(t *testing.T) {
	called := false
	service := NewService(
		runnerFunc(func(context.Context, domain.ActorContext, domain.AnalysisRequest) (domain.AnalysisResult, error) {
			called = true
			return domain.AnalysisResult{}, nil
		}),
		validatorFunc(func(context.Context, domain.AnalysisRequest, domain.ChartSpec) error { return nil }),
	)

	_, err := service.Analyze(context.Background(), actorWithDatasource("prometheus-main"), validRequest("other"))
	if domain.ErrorCodeOf(err) != domain.CodeDatasourceForbidden {
		t.Fatalf("expected DATASOURCE_FORBIDDEN, got %v", err)
	}
	if called {
		t.Fatal("agent must not run before access validation")
	}
}

func TestServicePropagatesAgentFailure(t *testing.T) {
	want := domain.NewError(domain.CodeMCPUnavailable, "MCP is unavailable")
	service := NewService(
		runnerFunc(func(context.Context, domain.ActorContext, domain.AnalysisRequest) (domain.AnalysisResult, error) {
			return domain.AnalysisResult{}, want
		}),
		validatorFunc(func(context.Context, domain.AnalysisRequest, domain.ChartSpec) error { return nil }),
	)

	_, err := service.Analyze(context.Background(), actorWithDatasource("prometheus-main"), validRequest("prometheus-main"))
	if !errors.Is(err, want) {
		t.Fatalf("expected agent error to propagate, got %v", err)
	}
}

func TestServiceRunsPolicyValidationAfterAgent(t *testing.T) {
	validated := false
	request := validRequest("prometheus-main")
	result := domain.AnalysisResult{Chart: domain.ChartSpec{DatasourceUID: request.Scope.DatasourceUID}}
	service := NewService(
		runnerFunc(func(context.Context, domain.ActorContext, domain.AnalysisRequest) (domain.AnalysisResult, error) {
			return result, nil
		}),
		validatorFunc(func(_ context.Context, gotRequest domain.AnalysisRequest, gotSpec domain.ChartSpec) error {
			validated = gotRequest.Scope.DatasourceUID == gotSpec.DatasourceUID
			return nil
		}),
	)

	if _, err := service.Analyze(context.Background(), actorWithDatasource("prometheus-main"), request); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if !validated {
		t.Fatal("ChartSpec validator was not called after the agent")
	}
}

func actorWithDatasource(uid string) domain.ActorContext {
	return domain.ActorContext{
		OrgID: 1,
		Login: "alice",
		Access: domain.AccessScope{
			AllowedDatasourceUIDs: []string{uid},
		},
	}
}

func validRequest(uid string) domain.AnalysisRequest {
	return domain.AnalysisRequest{
		Text: "checkout request rate",
		Scope: domain.AnalysisScope{
			DatasourceUID: uid,
			TimeRange: domain.TimeRange{
				From: "now-30m",
				To:   "now",
			},
		},
	}
}
