package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

const maxGrafanaQueryResponseBytes = 4 << 20

const maxPromQLBytes = 16 << 10

var (
	errDatasourceNotAllowed = errors.New("data source is not allowed")
	errInvalidGrafanaQuery  = errors.New("invalid Grafana query")
	errGrafanaUnavailable   = errors.New("Grafana query service unavailable")
	errGrafanaRejected      = errors.New("Grafana rejected the query")
	errGrafanaResponseLarge = errors.New("Grafana query response is too large")
)

type appJSONData struct {
	PrometheusDatasourceUIDs []string `json:"prometheusDatasourceUids"`
}

type prometheusRangeQuery struct {
	DatasourceUID string
	Expr          string
	From          time.Time
	To            time.Time
	Step          time.Duration
	MaxDataPoints int
}

type grafanaAPIClient struct {
	baseURL               *url.URL
	bearerToken           string
	allowedDatasourceUIDs map[string]struct{}
	httpClient            *http.Client
	timeout               time.Duration
}

type grafanaDatasourceRef struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

type grafanaPrometheusQuery struct {
	RefID         string               `json:"refId"`
	Expr          string               `json:"expr"`
	Format        string               `json:"format"`
	IntervalMS    int64                `json:"intervalMs"`
	MaxDataPoints int                  `json:"maxDataPoints"`
	Datasource    grafanaDatasourceRef `json:"datasource"`
}

type grafanaQueryRequest struct {
	Queries []grafanaPrometheusQuery `json:"queries"`
	From    string                   `json:"from"`
	To      string                   `json:"to"`
}

// newGrafanaAPIClientFromContext uses Grafana's managed app service account.
// Grafana injects the API URL and token into the SDK context; neither value is
// accepted from the browser or forwarded to the MCP server.
func newGrafanaAPIClientFromContext(
	ctx context.Context,
	httpClient *http.Client,
	timeout time.Duration,
) (*grafanaAPIClient, error) {
	pluginContext := backend.PluginConfigFromContext(ctx)
	if pluginContext.AppInstanceSettings == nil {
		return nil, fmt.Errorf("Grafana app settings: %w", errGrafanaUnavailable)
	}
	var settings appJSONData
	if err := json.Unmarshal(pluginContext.AppInstanceSettings.JSONData, &settings); err != nil {
		return nil, fmt.Errorf("Grafana app settings: %w", errGrafanaUnavailable)
	}

	cfg := backend.GrafanaConfigFromContext(ctx)
	appURL, err := cfg.AppURL()
	if err != nil {
		return nil, fmt.Errorf("Grafana app URL: %w", errGrafanaUnavailable)
	}
	token, err := cfg.PluginAppClientSecret()
	if err != nil || strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("Grafana managed service account: %w", errGrafanaUnavailable)
	}

	parsedURL, err := url.Parse(appURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return nil, fmt.Errorf("Grafana app URL: %w", errGrafanaUnavailable)
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	allowed := make(map[string]struct{}, len(settings.PrometheusDatasourceUIDs))
	for _, uid := range settings.PrometheusDatasourceUIDs {
		uid = strings.TrimSpace(uid)
		if uid != "" {
			allowed[uid] = struct{}{}
		}
	}

	return &grafanaAPIClient{
		baseURL:               parsedURL,
		bearerToken:           token,
		allowedDatasourceUIDs: allowed,
		httpClient:            httpClient,
		timeout:               timeout,
	}, nil
}

func (c *grafanaAPIClient) queryPrometheusRange(ctx context.Context, query prometheusRangeQuery) (json.RawMessage, error) {
	if _, ok := c.allowedDatasourceUIDs[query.DatasourceUID]; !ok {
		return nil, errDatasourceNotAllowed
	}
	if strings.TrimSpace(query.Expr) == "" || len(query.Expr) > maxPromQLBytes || query.From.IsZero() || query.To.IsZero() || !query.From.Before(query.To) {
		return nil, errInvalidGrafanaQuery
	}
	if query.Step < time.Millisecond {
		return nil, errInvalidGrafanaQuery
	}
	if query.MaxDataPoints <= 0 || query.MaxDataPoints > 11000 {
		return nil, errInvalidGrafanaQuery
	}

	payload := grafanaQueryRequest{
		From: strconv.FormatInt(query.From.UnixMilli(), 10),
		To:   strconv.FormatInt(query.To.UnixMilli(), 10),
		Queries: []grafanaPrometheusQuery{
			{
				RefID:         "A",
				Expr:          query.Expr,
				Format:        "time_series",
				IntervalMS:    query.Step.Milliseconds(),
				MaxDataPoints: query.MaxDataPoints,
				Datasource: grafanaDatasourceRef{
					Type: "prometheus",
					UID:  query.DatasourceUID,
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errInvalidGrafanaQuery
	}

	queryURL, err := url.JoinPath(c.baseURL.String(), "api", "ds", "query")
	if err != nil {
		return nil, fmt.Errorf("Grafana query URL: %w", errGrafanaUnavailable)
	}
	requestCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, queryURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("Grafana query request: %w", errGrafanaUnavailable)
	}
	req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Grafana query request: %w", errGrafanaUnavailable)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// Do not propagate an upstream response body: it can contain configuration
		// details and must never echo credentials to the frontend or LLM.
		return nil, fmt.Errorf("%w (status %d)", errGrafanaRejected, resp.StatusCode)
	}

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, maxGrafanaQueryResponseBytes+1))
	if err != nil {
		return nil, errGrafanaUnavailable
	}
	if len(responseBody) > maxGrafanaQueryResponseBytes {
		return nil, errGrafanaResponseLarge
	}
	if !json.Valid(responseBody) {
		return nil, errGrafanaUnavailable
	}

	return json.RawMessage(responseBody), nil
}
