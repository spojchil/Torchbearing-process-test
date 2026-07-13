# Torchbearing 当前已实现功能清单

> 更新日期：2026-07-13  
> 当前阶段：MS1  
> 交付形态：可编译、可测试、使用确定性 mock 串联的 MVP 骨架

## 1. 项目当前能力概览

Torchbearing 当前实现了一条最小指标分析流程：接收自然语言分析请求和 Grafana 范围信息，经过 Scope 校验、Mock Agent 计划生成、Mock Metrics 查询和图表定义构建，最终返回文字说明、图表契约或强类型错误。

当前所有数据均为 deterministic mock/in-memory 数据，不连接真实 Grafana、Prometheus、LLM、数据库或第三方服务。

## 2. 后端功能

| 编号 | 功能 | 当前状态 | 说明 | 主要代码位置 |
|---|---|---|---|---|
| BE-01 | 强类型分析契约 | 已实现 | 定义请求、响应、图表、时间范围和错误结构 | `contracts/`、`internal/contracts/`、`internal/core/` |
| BE-02 | Scope 规范化 | 已实现 | 清理 datasource UID、开始时间和结束时间的首尾空白 | `internal/scope/` |
| BE-03 | Scope 必填校验 | 已实现 | 校验 datasource UID、开始时间和结束时间非空 | `internal/scope/` |
| BE-04 | 时间范围校验 | 已实现 | 支持 RFC3339 和 `now±数字{s|m|h|d|w}` 的先后比较 | `internal/scope/` |
| BE-05 | Agent SDK abstraction | 已实现 | 使用强类型请求生成分析计划，不暴露具体实现 | `sdk/agent/` |
| BE-06 | Mock Agent | 已实现（Mock） | 根据固定场景生成确定性 PromQL、消息和图表意图 | `mocks/agent/` |
| BE-07 | Metrics SDK abstraction | 已实现 | 使用强类型 MetricQuery/MetricResult 隔离数据源实现 | `sdk/metrics/` |
| BE-08 | Mock Metrics | 已实现（Mock） | 返回固定三点序列、空结果、固定失败或单点边界数据 | `mocks/metrics/` |
| BE-09 | ChartBuilder | 已实现 | 将 AnalysisPlan 和 MetricResult 转换为 renderer-neutral ChartSpec | `internal/chart/` |
| BE-10 | Typed Error | 已实现 | 固定错误码、message、retryable、request ID 和错误链 | `internal/core/errors.go` |
| BE-11 | 确定性 Request ID | 已实现（Mock） | 生成 `mock-analysis-001` 等稳定 ID | `mocks/deterministic/` |
| BE-12 | 固定 Clock | 已实现（Mock） | 返回注入的固定时间，不读取系统时钟 | `mocks/deterministic/` |
| BE-13 | 分析主流程编排 | 已实现 | 串联 Scope、Agent、Metrics 和 ChartBuilder | `internal/analysis/` |
| BE-14 | 失败短路 | 已实现 | 任一模块失败后停止下游调用，并补充 request ID | `internal/analysis/` |
| BE-15 | 内存 Gateway | 已实现（In-memory） | 在传输 DTO 与领域模型之间转换，不启动 HTTP/RPC | `internal/transport/grafana/` |
| BE-16 | 依赖组装 | 已实现 | 按场景注入 B/C 的确定性实现 | `internal/bootstrap/` |
| BE-17 | CLI Demo | 已实现 | 从固定请求运行完整成功流程并输出 JSON | `cmd/torchbearing/` |

## 3. 前端功能

| 编号 | 功能 | 当前状态 | 说明 | 主要代码位置 |
|---|---|---|---|---|
| FE-01 | TypeScript 分析契约 | 已实现 | 定义 AnalysisRequest、AnalysisResponse、ChartSpec 和 ErrorResponse | `plugin/src/contracts/` |
| FE-02 | AnalysisSDK | 已实现 | 页面只通过业务 SDK 发起分析，不直接访问内部实现 | `plugin/src/sdk/` |
| FE-03 | Fixture Registry | 已实现（Mock） | 加载公共 JSON fixtures，并转换为强类型前端对象 | `plugin/src/mocks/fixtureRegistry.ts` |
| FE-04 | DeterministicAnalysisSDK | 已实现（Mock） | 稳定复现成功、空结果、Agent 失败、Metrics 失败和无效范围 | `plugin/src/mocks/DeterministicAnalysisSDK.ts` |
| FE-05 | 默认 Scope | 已实现 | 默认使用 `prometheus-mock` 和 `now-30m → now` | `plugin/src/features/scope/` |
| FE-06 | Scope 表单模型 | 已实现 | 输入清理、必填校验、SDK Scope 转换和摘要生成 | `plugin/src/features/scope/` |
| FE-07 | Analysis Workbench | 已实现 | 管理分析请求和展示状态 | `plugin/src/features/analysis/` |
| FE-08 | Workbench 状态机 | 已实现 | 支持 `idle/loading/success/empty/error` | `plugin/src/features/analysis/AnalysisWorkbench.ts` |
| FE-09 | Typed Error 展示模型 | 已实现 | 将 AnalysisSDKError 转换为稳定 error 状态 | `plugin/src/features/analysis/` |
| FE-10 | 最小 Grafana 页面入口 | 已实现（骨架） | 提供输入框、分析按钮、状态、图表标题和 PromQL 展示 | `plugin/src/module.tsx` |

## 4. 已实现的 SDK/接口

当前实现了 8 个主要功能接口：

| 接口 | 实现 |
|---|---|
| `core.ScopeResolver` | `scope.Resolver` |
| `core.ChartBuilder` | `chart.Builder` |
| `core.IDGenerator` | `deterministic.IDGenerator` |
| `core.Clock` | `deterministic.Clock` |
| `core.Analyzer` | `analysis.Service` |
| `sdk/agent.Client` | `mocks/agent.Client` |
| `sdk/metrics.Client` | `mocks/metrics.Client` |
| `AnalysisSDK` | `DeterministicAnalysisSDK` |

此外，后端提供 `Gateway.Analyze` 内存调用入口和 `cmd/torchbearing` CLI Demo 入口。目前没有真实 HTTP/RPC endpoint。

## 5. Mock 场景

| 场景 | Request ID | 结果 |
|---|---|---|
| 成功 | `mock-analysis-001` | 固定 checkout PromQL、三点数据、一个 timeseries 图表 |
| 空结果 | `mock-analysis-002` | `DataStateEmpty`、空图表数组 |
| Agent 失败 | `mock-analysis-003` | retryable `AGENT_UNAVAILABLE` |
| Metrics 失败 | `mock-analysis-004` | retryable `METRICS_UNAVAILABLE` |
| 单点边界 | `mock-analysis-005` | 单点零值数据、一个 stat 图表 |
| 无效时间范围 | `mock-analysis-005` | `INVALID_SCOPE`，不调用 Agent/Metrics |
| 空分析文本 | 固定边界 ID | `INVALID_ARGUMENT` |
| 请求取消 | 当前请求 ID | 保留 context cancellation 错误链 |

所有场景都不使用随机值、真实当前时间或外部数据。

## 6. 当前主流程

```text
AnalysisRequest
  → ScopeResolver
  → Agent SDK / Mock Agent
  → Metrics SDK / Mock Metrics
  → ChartBuilder
  → AnalysisResponse 或 Typed Error
```

前端流程：

```text
Grafana 最小页面
  → AnalysisWorkbench
  → AnalysisSDK
  → Fixture Mock
  → success / empty / error 展示状态
```

## 7. 测试功能

当前已经添加：

- JSON Schema 和公共 fixture 契约测试。
- Go SDK 与 TypeScript SDK 契约检查。
- Typed error 单元测试。
- Scope 成功、空输入、失败、取消和边界测试。
- Agent mock 场景测试。
- Metrics mock 场景测试。
- Clock 和 ID Generator 确定性测试。
- ChartBuilder 成功、空数据、失败和边界测试。
- 前端 AnalysisSDK、Scope Model 和 Workbench 测试。
- CLI 输出测试。
- 完整 MS1 后端集成测试。
- Go race test。

## 8. 如何验证

### 8.1 运行 Demo

在项目根目录执行：

```bash
go run ./cmd/torchbearing
```

预期返回 `mock-analysis-001` 和一个固定 timeseries 图表定义。

### 8.2 运行主流程集成测试

```bash
go test -v ./tests/integration
```

### 8.3 运行完整测试

```bash
go test -v ./...
go vet ./...
go build ./...
```

### 8.4 前端测试

当前缺少可用的 `plugin/package-lock.json` 和 `plugin/node_modules`。依赖完成锁定和安装后执行：

```bash
cd plugin
npm run typecheck
npm run test:ci
npm run lint
npm run format:check
npm run build
```

## 9. 尚未实现

以下能力当前没有实现：

- 真实 Grafana 数据源查询。
- 真实 Prometheus/PromQL 网络请求。
- 真实 LLM 或 Eino Agent。
- MCP client/server。
- HTTP、SSE、WebSocket 或 RPC transport。
- 数据库存储和会话持久化。
- 多轮对话和流式输出。
- 多图表画布编辑与持久化。
- Skills、知识库、Playbook、告警分析、审计和发布流程。
- 真实部署与生产配置。

## 10. 已知集成问题

1. `ChartSpec` 尚未携带指标点位，前端只能展示图表元信息和 PromQL。
2. `plugin.json` 声明的 backend executable 与当前 Go CLI 入口尚未对齐。
3. 前端依赖和 lockfile 尚未完成，无法执行正式 webpack/Jest/ESLint/Prettier 验收。
4. 当前前后端通过共享契约和 fixtures 对齐，没有建立真实网络连接。

## 11. 相关文档

- [MS1 完整实现说明](./MS1_IMPLEMENTATION.md)
- [软件工程规范](../exclude/软件工程规范.md)
- [架构设计规范](../exclude/架构设计规范.md)
- [产品与架构 Proposals](../exclude/proposals/)
