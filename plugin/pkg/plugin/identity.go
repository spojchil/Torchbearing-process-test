package plugin

import (
	"context"
	"errors"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

var errUnauthenticated = errors.New("authenticated Grafana user is required")

// callerIdentity is the deliberately minimized identity that may cross the
// Plugin Backend -> MCP boundary. Email, cookies and authorization headers are
// intentionally excluded.
type callerIdentity struct {
	OrgID int64
	Login string
	Role  string
}

func callerIdentityFromContext(ctx context.Context) (callerIdentity, error) {
	pluginContext := backend.PluginConfigFromContext(ctx)
	if pluginContext.OrgID <= 0 || pluginContext.User == nil || strings.TrimSpace(pluginContext.User.Login) == "" {
		return callerIdentity{}, errUnauthenticated
	}

	return callerIdentity{
		OrgID: pluginContext.OrgID,
		Login: pluginContext.User.Login,
		Role:  pluginContext.User.Role,
	}, nil
}
