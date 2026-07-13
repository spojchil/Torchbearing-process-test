package plugin

import (
	"context"
	"errors"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestCallerIdentityFromContextMinimizesIdentity(t *testing.T) {
	ctx := backend.WithPluginContext(context.Background(), backend.PluginContext{
		OrgID: 7,
		User: &backend.User{
			Login: "alice",
			Name:  "Alice Example",
			Email: "alice@example.com",
			Role:  "Viewer",
		},
		AppInstanceSettings: &backend.AppInstanceSettings{
			JSONData: []byte(`{"prometheusDatasourceUids":[" prometheus-main ","prometheus-main"]}`),
		},
	})

	got, err := callerIdentityFromContext(ctx)
	if err != nil {
		t.Fatalf("caller identity: %v", err)
	}
	if got.OrgID != 7 || got.Login != "alice" || got.Role != "Viewer" {
		t.Fatalf("unexpected minimized identity: %#v", got)
	}
	if len(got.Access.AllowedDatasourceUIDs) != 1 || got.Access.AllowedDatasourceUIDs[0] != "prometheus-main" {
		t.Fatalf("unexpected access scope: %#v", got.Access)
	}
}

func TestCallerIdentityFromContextRejectsInvalidSettings(t *testing.T) {
	ctx := backend.WithPluginContext(context.Background(), backend.PluginContext{
		OrgID:               7,
		User:                &backend.User{Login: "alice"},
		AppInstanceSettings: &backend.AppInstanceSettings{JSONData: []byte(`{"broken"`)},
	})

	_, err := callerIdentityFromContext(ctx)
	if !errors.Is(err, errInvalidAppSettings) {
		t.Fatalf("expected errInvalidAppSettings, got %v", err)
	}
}

func TestCallerIdentityFromContextRequiresSDKIdentity(t *testing.T) {
	_, err := callerIdentityFromContext(context.Background())
	if !errors.Is(err, errUnauthenticated) {
		t.Fatalf("expected errUnauthenticated, got %v", err)
	}
}
