# Torchbearing Grafana App Plugin

MS1 空骨架包含：

- 一个“指标分析工作台”页面；
- 一个 TypeScript API/Model 层；
- 一个 Go Backend Resource Handler；
- `POST /analysis` Mock 接口；
- Grafana Plugin SDK 健康检查；
- Jest、Go test和Playwright测试骨架。

`/analysis` 当前明确返回 Mock ChartSpec，不查询真实 Prometheus。真实指标检索、MCP和查询执行分别由后续Issue实现。

## 本地验证

```bash
npm ci
npm run typecheck
npm run lint
npm run test:ci
npm run build
go test ./pkg/...
go build ./pkg/...
```

使用 `npm run server` 启动开发环境。修改 `src/plugin.json` 后必须重启 Grafana 开发服务器。

## 官方脚手架说明

该目录由 `@grafana/create-plugin` 7.8.1 生成并按 Issue #4 收敛。构建配置保留在受工具管理的 `.config/` 中。

- [App Plugin开发指南](https://grafana.com/developers/plugin-tools/how-to-guides/app-plugins/)
- [Plugin Backend说明](https://grafana.com/developers/plugin-tools/key-concepts/backend-plugins/)
- [plugin.json参考](https://grafana.com/developers/plugin-tools/reference/plugin-json/)
