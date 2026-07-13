package grafana

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

const testManagedServiceAccountToken = "test-managed-service-account-token"

func TestGrafanaAPIClientQueriesAllowedPrometheusDatasource(t *testing.T) {
	var called atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		called.Add(1)
		if req.Method != http.MethodPost || req.URL.Path != "/api/ds/query" {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer "+testManagedServiceAccountToken {
			t.Errorf("unexpected authorization header: %q", got)
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read request: %v", err)
		}
		var payload struct {
			Queries []struct {
				Expr       string `json:"expr"`
				Datasource struct {
					Type string `json:"type"`
					UID  string `json:"uid"`
				} `json:"datasource"`
			} `json:"queries"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if len(payload.Queries) != 1 || payload.Queries[0].Datasource.UID != "prometheus-main" || payload.Queries[0].Datasource.Type != "prometheus" {
			t.Fatalf("unexpected query payload: %#v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"results":{"A":{"frames":[]}}}`)
	}))
	defer server.Close()

	client := newTestGrafanaAPIClient(t, server.URL, []string{"prometheus-main"}, server.Client(), time.Second)
	response, err := client.queryPrometheusRange(context.Background(), validPrometheusRangeQuery("prometheus-main"))
	if err != nil {
		t.Fatalf("query Prometheus range: %v", err)
	}
	if !json.Valid(response) || called.Load() != 1 {
		t.Fatalf("expected one valid JSON response, got %q and %d calls", response, called.Load())
	}
}

func TestGrafanaAPIClientRejectsDatasourceBeforeNetwork(t *testing.T) {
	var called atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		called.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestGrafanaAPIClient(t, server.URL, []string{"prometheus-main"}, server.Client(), time.Second)
	_, err := client.queryPrometheusRange(context.Background(), validPrometheusRangeQuery("other-datasource"))
	if !errors.Is(err, errDatasourceNotAllowed) {
		t.Fatalf("expected errDatasourceNotAllowed, got %v", err)
	}
	if called.Load() != 0 {
		t.Fatalf("disallowed datasource reached Grafana %d times", called.Load())
	}
}

func TestGrafanaAPIClientRedactsRejectedResponse(t *testing.T) {
	const upstreamSecret = "upstream-sensitive-detail"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, upstreamSecret, http.StatusForbidden)
	}))
	defer server.Close()

	client := newTestGrafanaAPIClient(t, server.URL, []string{"prometheus-main"}, server.Client(), time.Second)
	_, err := client.queryPrometheusRange(context.Background(), validPrometheusRangeQuery("prometheus-main"))
	if !errors.Is(err, errGrafanaRejected) {
		t.Fatalf("expected errGrafanaRejected, got %v", err)
	}
	if strings.Contains(err.Error(), upstreamSecret) || strings.Contains(err.Error(), testManagedServiceAccountToken) {
		t.Fatalf("error leaked sensitive data: %v", err)
	}
}

func TestGrafanaAPIClientEnforcesTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"results":{}}`)
	}))
	defer server.Close()

	client := newTestGrafanaAPIClient(t, server.URL, []string{"prometheus-main"}, server.Client(), 20*time.Millisecond)
	_, err := client.queryPrometheusRange(context.Background(), validPrometheusRangeQuery("prometheus-main"))
	if !errors.Is(err, errGrafanaUnavailable) {
		t.Fatalf("expected errGrafanaUnavailable, got %v", err)
	}
}

func newTestGrafanaAPIClient(
	t *testing.T,
	appURL string,
	allowedDatasourceUIDs []string,
	httpClient *http.Client,
	timeout time.Duration,
) *grafanaAPIClient {
	t.Helper()
	ctx := backend.WithGrafanaConfig(context.Background(), backend.NewGrafanaCfg(map[string]string{
		backend.AppURL:          appURL,
		backend.AppClientSecret: testManagedServiceAccountToken,
	}))
	settings, err := json.Marshal(appJSONData{PrometheusDatasourceUIDs: allowedDatasourceUIDs})
	if err != nil {
		t.Fatalf("marshal app settings: %v", err)
	}
	ctx = backend.WithPluginContext(ctx, backend.PluginContext{
		AppInstanceSettings: &backend.AppInstanceSettings{JSONData: settings},
	})
	client, err := newGrafanaAPIClientFromContext(ctx, httpClient, timeout)
	if err != nil {
		t.Fatalf("new Grafana API client: %v", err)
	}
	return client
}

func validPrometheusRangeQuery(datasourceUID string) prometheusRangeQuery {
	to := time.Unix(1800000000, 0)
	return prometheusRangeQuery{
		DatasourceUID: datasourceUID,
		Expr:          `sum(rate(http_requests_total[5m]))`,
		From:          to.Add(-30 * time.Minute),
		To:            to,
		Step:          15 * time.Second,
		MaxDataPoints: 120,
	}
}
