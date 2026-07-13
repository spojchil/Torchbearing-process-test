// Package grafana 提供 MS1 的强类型内存网关，不启动 HTTP/RPC 服务。
package grafana

import (
	"context"
	"errors"

	"github.com/spojchil/torchbearing/internal/contracts"
	"github.com/spojchil/torchbearing/internal/core"
)

// Gateway 将稳定传输契约映射到核心 Analyzer 接口。
type Gateway struct {
	analyzer core.Analyzer
}

// NewGateway 创建不连接真实 Grafana 服务的内存网关。
func NewGateway(analyzer core.Analyzer) *Gateway {
	return &Gateway{analyzer: analyzer}
}

// Analyze 返回成功响应或强类型错误响应，二者严格二选一。
func (g *Gateway) Analyze(
	ctx context.Context,
	request contracts.AnalysisRequest,
) (contracts.AnalysisResponse, *contracts.ErrorResponse) {
	response, err := g.analyzer.Analyze(ctx, core.AnalysisRequest{
		Text: request.Text,
		Scope: core.AnalysisScope{
			DatasourceUID: request.Scope.DatasourceUID,
			TimeRange: core.TimeRange{
				From: request.Scope.TimeRange.From,
				To:   request.Scope.TimeRange.To,
			},
		},
	})
	if err != nil {
		return contracts.AnalysisResponse{}, toErrorResponse(err)
	}

	charts := make([]contracts.ChartSpec, 0, len(response.Charts))
	for _, chart := range response.Charts {
		charts = append(charts, contracts.ChartSpec{
			ID:            chart.ID,
			Title:         chart.Title,
			Type:          chart.Type,
			DatasourceUID: chart.DatasourceUID,
			PromQL:        chart.PromQL,
			TimeRange: contracts.TimeRange{
				From: chart.TimeRange.From,
				To:   chart.TimeRange.To,
			},
		})
	}

	return contracts.AnalysisResponse{
		RequestID: response.RequestID,
		Message:   response.Message,
		Charts:    charts,
		Mock:      response.Mock,
	}, nil
}

func toErrorResponse(err error) *contracts.ErrorResponse {
	var typed *core.AppError
	if errors.As(err, &typed) {
		return &contracts.ErrorResponse{
			Code:      typed.Code,
			Message:   typed.Message,
			Retryable: typed.Retryable,
			RequestID: typed.RequestID,
		}
	}
	return &contracts.ErrorResponse{
		Code:      core.ErrorCodeInternal,
		Message:   "analysis failed unexpectedly",
		Retryable: false,
	}
}
