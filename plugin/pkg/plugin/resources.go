package plugin

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type analysisRequest struct {
	Text  string        `json:"text"`
	Scope analysisScope `json:"scope"`
}

type analysisScope struct {
	DatasourceUID string    `json:"datasourceUid"`
	TimeRange     timeRange `json:"timeRange"`
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
}

type analysisResponse struct {
	RequestID string    `json:"requestId"`
	Chart     chartSpec `json:"chart"`
	Mock      bool      `json:"mock"`
}

func (a *App) handleAnalysis(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if _, err := callerIdentityFromContext(req.Context()); err != nil {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}

	var body analysisRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(body.Text) == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}
	if body.Scope.TimeRange.From == "" {
		body.Scope.TimeRange.From = "now-30m"
	}
	if body.Scope.TimeRange.To == "" {
		body.Scope.TimeRange.To = "now"
	}

	response := analysisResponse{
		RequestID: "mock-" + time.Now().UTC().Format("20060102T150405.000000000Z"),
		Mock:      true,
		Chart: chartSpec{
			ID:            "request-rate",
			Title:         "Checkout 请求速率（Mock）",
			Type:          "timeseries",
			DatasourceUID: body.Scope.DatasourceUID,
			PromQL:        `sum(rate(http_requests_total{service="checkout"}[5m]))`,
			TimeRange:     body.Scope.TimeRange,
		},
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/analysis", a.handleAnalysis)
}
