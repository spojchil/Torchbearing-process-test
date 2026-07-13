package plugin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
)

const maxAnalysisRequestBytes = 64 << 10

type analysisRequest struct {
	Text  string        `json:"text"`
	Scope analysisScope `json:"scope"`
}

type analysisScope struct {
	DatasourceUID string    `json:"datasourceUid"`
	TimeRange     timeRange `json:"timeRange"`
	Service       string    `json:"service,omitempty"`
	Environment   string    `json:"environment,omitempty"`
	Namespace     string    `json:"namespace,omitempty"`
	Cluster       string    `json:"cluster,omitempty"`
}

type timeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type chartSpec struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Type          string    `json:"type"`
	DatasourceUID string    `json:"datasourceUid"`
	PromQL        string    `json:"promql"`
	TimeRange     timeRange `json:"timeRange"`
	Unit          string    `json:"unit,omitempty"`
	Legend        string    `json:"legend,omitempty"`
}

type querySummary struct {
	Language    string `json:"language"`
	Expression  string `json:"expression"`
	Status      string `json:"status"`
	DurationMS  int64  `json:"durationMs"`
	SeriesCount int    `json:"seriesCount"`
}

type evidence struct {
	Metrics     []string `json:"metrics"`
	Explanation string   `json:"explanation"`
}

type analysisResponse struct {
	RequestID string       `json:"requestId"`
	Chart     chartSpec    `json:"chart"`
	Query     querySummary `json:"query"`
	Evidence  evidence     `json:"evidence"`
	Mock      bool         `json:"mock"`
}

type errorEnvelope struct {
	RequestID string        `json:"requestId"`
	Error     responseError `json:"error"`
}

type responseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (a *App) handleAnalysis(w http.ResponseWriter, req *http.Request) {
	requestID := "req-" + time.Now().UTC().Format("20060102T150405.000000000Z")
	if req.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, requestID, "INVALID_REQUEST", "method not allowed")
		return
	}

	actor, err := callerIdentityFromContext(req.Context())
	if err != nil {
		code := "UNAUTHENTICATED"
		message := "authentication required"
		status := http.StatusUnauthorized
		if errors.Is(err, errInvalidAppSettings) {
			code = "INTERNAL"
			message = "plugin configuration is invalid"
			status = http.StatusInternalServerError
		}
		writeJSONError(w, status, requestID, code, message)
		return
	}

	req.Body = http.MaxBytesReader(w, req.Body, maxAnalysisRequestBytes)
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	var body analysisRequest
	if err := decoder.Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, requestID, "INVALID_REQUEST", "invalid request body")
		return
	}
	body.Text = strings.TrimSpace(body.Text)
	if body.Scope.TimeRange.From == "" {
		body.Scope.TimeRange.From = "now-30m"
	}
	if body.Scope.TimeRange.To == "" {
		body.Scope.TimeRange.To = "now"
	}

	result, err := a.analysisService.Analyze(req.Context(), actor, toDomainRequest(body))
	if err != nil {
		writeDomainError(w, requestID, err)
		return
	}

	writeJSON(w, http.StatusOK, toAnalysisResponse(requestID, result))
}

func toDomainRequest(request analysisRequest) domain.AnalysisRequest {
	return domain.AnalysisRequest{
		Text: request.Text,
		Scope: domain.AnalysisScope{
			DatasourceUID: request.Scope.DatasourceUID,
			TimeRange: domain.TimeRange{
				From: request.Scope.TimeRange.From,
				To:   request.Scope.TimeRange.To,
			},
			Service:     request.Scope.Service,
			Environment: request.Scope.Environment,
			Namespace:   request.Scope.Namespace,
			Cluster:     request.Scope.Cluster,
		},
	}
}

func toAnalysisResponse(requestID string, result domain.AnalysisResult) analysisResponse {
	return analysisResponse{
		RequestID: requestID,
		Mock:      true,
		Chart: chartSpec{
			ID:            result.Chart.ID,
			Title:         result.Chart.Title,
			Type:          string(result.Chart.Type),
			DatasourceUID: result.Chart.DatasourceUID,
			PromQL:        result.Chart.PromQL,
			TimeRange: timeRange{
				From: result.Chart.TimeRange.From,
				To:   result.Chart.TimeRange.To,
			},
			Unit:   result.Chart.Unit,
			Legend: result.Chart.Legend,
		},
		Query: querySummary{
			Language:    result.Query.Language,
			Expression:  result.Query.Expression,
			Status:      string(result.Query.Status),
			DurationMS:  result.Query.Duration.Milliseconds(),
			SeriesCount: result.Query.SeriesCount,
		},
		Evidence: evidence{
			Metrics:     result.Evidence.Metrics,
			Explanation: result.Evidence.Explanation,
		},
	}
}

func writeDomainError(w http.ResponseWriter, requestID string, err error) {
	code := domain.ErrorCodeOf(err)
	status := http.StatusInternalServerError
	switch code {
	case domain.CodeInvalidRequest, domain.CodeIntentAmbiguous, domain.CodeQueryInvalid:
		status = http.StatusBadRequest
	case domain.CodeMetricNotFound:
		status = http.StatusNotFound
	case domain.CodeDatasourceForbidden:
		status = http.StatusForbidden
	case domain.CodeNoData, domain.CodeQueryLimited:
		status = http.StatusUnprocessableEntity
	case domain.CodeModelUnavailable, domain.CodeMCPUnavailable:
		status = http.StatusServiceUnavailable
	}
	writeJSONError(w, status, requestID, string(code), domain.SafeMessage(err))
}

func writeJSONError(w http.ResponseWriter, status int, requestID, code, message string) {
	writeJSON(w, status, errorEnvelope{
		RequestID: requestID,
		Error: responseError{
			Code:    code,
			Message: message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (a *App) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/analysis", a.handleAnalysis)
}
