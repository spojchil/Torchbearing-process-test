package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// mockCallResourceResponseSender implements backend.CallResourceResponseSender
// for use in tests.
type mockCallResourceResponseSender struct {
	response *backend.CallResourceResponse
}

// Send sets the received *backend.CallResourceResponse to s.response
func (s *mockCallResourceResponseSender) Send(response *backend.CallResourceResponse) error {
	s.response = response
	return nil
}

// TestCallResource tests CallResource calls, using backend.CallResourceRequest and backend.CallResourceResponse.
// This ensures the httpadapter for CallResource works correctly.
func TestCallResource(t *testing.T) {
	// Initialize app
	inst, err := NewApp(context.Background(), backend.AppInstanceSettings{})
	if err != nil {
		t.Fatalf("new app: %s", err)
	}
	if inst == nil {
		t.Fatal("inst must not be nil")
	}
	app, ok := inst.(*App)
	if !ok {
		t.Fatal("inst must be of type *App")
	}

	// Set up and run test cases
	for _, tc := range []struct {
		name string

		method        string
		path          string
		body          []byte
		authenticated bool
		headers       map[string][]string

		expStatus int
		expBody   []byte
		expCode   string
	}{
		{
			name:      "get analysis 405",
			method:    http.MethodGet,
			path:      "analysis",
			expStatus: http.StatusMethodNotAllowed,
		},
		{
			name:          "post analysis 200",
			method:        http.MethodPost,
			path:          "analysis",
			body:          []byte(`{"text":"checkout request rate","scope":{"datasourceUid":"prometheus","timeRange":{"from":"now-30m","to":"now"}}}`),
			authenticated: true,
			expStatus:     http.StatusOK,
		},
		{
			name:          "post analysis requires text",
			method:        http.MethodPost,
			path:          "analysis",
			body:          []byte(`{"text":""}`),
			authenticated: true,
			expStatus:     http.StatusBadRequest,
			expCode:       "INVALID_REQUEST",
		},
		{
			name:          "post analysis rejects datasource outside access scope",
			method:        http.MethodPost,
			path:          "analysis",
			body:          []byte(`{"text":"checkout request rate","scope":{"datasourceUid":"other","timeRange":{"from":"now-30m","to":"now"}}}`),
			authenticated: true,
			expStatus:     http.StatusForbidden,
			expCode:       "DATASOURCE_FORBIDDEN",
		},
		{
			name:   "post analysis rejects spoofed identity header",
			method: http.MethodPost,
			path:   "analysis",
			body:   []byte(`{"text":"checkout request rate"}`),
			headers: map[string][]string{
				"X-Grafana-User": {"spoofed-admin"},
			},
			expStatus: http.StatusUnauthorized,
			expCode:   "UNAUTHENTICATED",
		},
		{
			name:      "get non existing handler 404",
			method:    http.MethodGet,
			path:      "not_found",
			expStatus: http.StatusNotFound,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Request by calling CallResource. This tests the httpadapter.
			var r mockCallResourceResponseSender
			pluginContext := backend.PluginContext{}
			if tc.authenticated {
				pluginContext = backend.PluginContext{
					OrgID: 42,
					User: &backend.User{
						Login: "alice",
						Name:  "Alice Example",
						Email: "alice@example.com",
						Role:  "Viewer",
					},
					AppInstanceSettings: &backend.AppInstanceSettings{
						JSONData: []byte(`{"prometheusDatasourceUids":["prometheus"]}`),
					},
				}
			}
			err = app.CallResource(context.Background(), &backend.CallResourceRequest{
				PluginContext: pluginContext,
				Method:        tc.method,
				Path:          tc.path,
				Headers:       tc.headers,
				Body:          tc.body,
			}, &r)
			if err != nil {
				t.Fatalf("CallResource error: %s", err)
			}
			if r.response == nil {
				t.Fatal("no response received from CallResource")
			}
			if tc.expStatus > 0 && tc.expStatus != r.response.Status {
				t.Errorf("response status should be %d, got %d", tc.expStatus, r.response.Status)
			}
			if len(tc.expBody) > 0 {
				if tb := bytes.TrimSpace(r.response.Body); !bytes.Equal(tb, tc.expBody) {
					t.Errorf("response body should be %s, got %s", tc.expBody, tb)
				}
			}
			if tc.expCode != "" {
				var body errorEnvelope
				if err := json.Unmarshal(r.response.Body, &body); err != nil {
					t.Fatalf("decode error response: %s", err)
				}
				if body.Error.Code != tc.expCode || body.RequestID == "" {
					t.Fatalf("expected error code %s with request ID, got %#v", tc.expCode, body)
				}
			}
			if tc.name == "post analysis 200" {
				var body analysisResponse
				if err := json.Unmarshal(r.response.Body, &body); err != nil {
					t.Fatalf("decode response: %s", err)
				}
				if !body.Mock || body.Chart.PromQL == "" || body.Query.Status != "success" {
					t.Fatalf("expected mock chart response, got %#v", body)
				}
				if len(body.Evidence.Metrics) != 1 || body.Evidence.Metrics[0] != "http_requests_total" {
					t.Fatalf("expected architecture path evidence, got %#v", body.Evidence)
				}
			}
		})
	}
}
