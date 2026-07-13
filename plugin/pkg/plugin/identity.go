package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	domain "github.com/spojchil/torchbearing/pkg/domain/analysis"
)

var errUnauthenticated = errors.New("authenticated Grafana user is required")

var errInvalidAppSettings = errors.New("invalid Grafana app settings")

type appSettings struct {
	PrometheusDatasourceUIDs []string `json:"prometheusDatasourceUids"`
}

// callerIdentityFromContext builds the deliberately minimized identity and
// trusted access scope that may cross the Plugin Backend -> MCP boundary.
// Email, cookies, and authorization headers are intentionally excluded.
func callerIdentityFromContext(ctx context.Context) (domain.ActorContext, error) {
	pluginContext := backend.PluginConfigFromContext(ctx)
	if pluginContext.OrgID <= 0 || pluginContext.User == nil || strings.TrimSpace(pluginContext.User.Login) == "" {
		return domain.ActorContext{}, errUnauthenticated
	}

	var settings appSettings
	if pluginContext.AppInstanceSettings != nil && len(pluginContext.AppInstanceSettings.JSONData) > 0 {
		if err := json.Unmarshal(pluginContext.AppInstanceSettings.JSONData, &settings); err != nil {
			return domain.ActorContext{}, errInvalidAppSettings
		}
	}

	allowed := make([]string, 0, len(settings.PrometheusDatasourceUIDs))
	seen := make(map[string]struct{}, len(settings.PrometheusDatasourceUIDs))
	for _, uid := range settings.PrometheusDatasourceUIDs {
		uid = strings.TrimSpace(uid)
		if uid == "" {
			continue
		}
		if _, exists := seen[uid]; exists {
			continue
		}
		seen[uid] = struct{}{}
		allowed = append(allowed, uid)
	}

	return domain.ActorContext{
		OrgID: pluginContext.OrgID,
		Login: pluginContext.User.Login,
		Role:  pluginContext.User.Role,
		Access: domain.AccessScope{
			AllowedDatasourceUIDs: allowed,
		},
	}, nil
}
