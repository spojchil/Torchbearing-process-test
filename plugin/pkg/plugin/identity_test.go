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
	})

	got, err := callerIdentityFromContext(ctx)
	if err != nil {
		t.Fatalf("caller identity: %v", err)
	}
	if got.OrgID != 7 || got.Login != "alice" || got.Role != "Viewer" {
		t.Fatalf("unexpected minimized identity: %#v", got)
	}
}

func TestCallerIdentityFromContextRequiresSDKIdentity(t *testing.T) {
	_, err := callerIdentityFromContext(context.Background())
	if !errors.Is(err, errUnauthenticated) {
		t.Fatalf("expected errUnauthenticated, got %v", err)
	}
}
