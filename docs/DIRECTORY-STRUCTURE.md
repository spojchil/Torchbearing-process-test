# Torchbearing 目录结构说明

本文单独说明仓库目录的职责、内容边界和维护人。具体产品范围、团队分工、接口字段和实现方案以对应 Proposal、Design Issue 为准。

## 1. 目录设计原则

- `cmd` 只负责进程入口和依赖装配，不放业务逻辑。
- `internal` 按领域拆分，领域规则留在领域目录中。
- `internal/platform` 只放可复用的技术适配，不承载业务概念。
- `api` 保存跨进程、跨语言契约，不保存服务实现。
- 四个 MCP Server 拥有独立入口和部署目录，但可以复用公共 MCP 库。
- 接口由能力提供方维护，消费者不复制领域模型。
- MS1 使用单一 Go module，不创建多个 module 或 `go.work`。
- 空目录通过 `.gitkeep` 保留，后续实现时可以删除对应占位文件。

## 2. 总目录树

```text
Torchbearing-process-test/
├── .github/
│   ├── ISSUE_TEMPLATE/
│   └── PULL_REQUEST_TEMPLATE.md
├── api/
│   ├── events/
│   ├── http/
│   └── mcp/
├── cmd/
│   ├── ai-core/
│   ├── grafana-app-backend/
│   ├── mcp-grafana/
│   ├── mcp-knowledge/
│   ├── mcp-playbook/
│   └── mcp-skills/
├── config/
├── data/
│   ├── services/
│   ├── skills/
│   │   ├── private/
│   │   └── shared/
│   └── playbooks/
│       ├── private/
│       └── shared/
├── deploy/
│   ├── docker/
│   │   ├── ai-core/
│   │   ├── mcp-grafana/
│   │   ├── mcp-knowledge/
│   │   ├── mcp-playbook/
│   │   └── mcp-skills/
│   └── k8s/
├── docs/
│   ├── DIRECTORY-STRUCTURE.md
│   └── PROCESS.md
├── internal/
│   ├── agent/
│   │   ├── application/
│   │   ├── hitl/
│   │   ├── orchestration/
│   │   └── transport/
│   ├── alerting/
│   │   ├── dispatcher/
│   │   └── receiver/
│   ├── audit/
│   ├── grafana/
│   │   └── mcpserver/
│   ├── grafanaplugin/
│   │   ├── identity/
│   │   └── proxy/
│   ├── knowledge/
│   │   ├── catalog/
│   │   ├── importer/
│   │   ├── mcpserver/
│   │   ├── retrieval/
│   │   └── store/sqlite/migrations/
│   ├── platform/
│   │   ├── config/
│   │   ├── grafana/
│   │   ├── mcp/
│   │   └── storage/
│   ├── playbook/
│   │   ├── engine/
│   │   ├── mcpserver/
│   │   └── repository/
│   ├── promotion/
│   │   ├── approval/
│   │   └── store/sqlite/migrations/
│   ├── session/
│   │   ├── canvas/
│   │   └── store/sqlite/migrations/
│   └── skills/
│       ├── mcpserver/
│       ├── middleware/
│       └── repository/
├── plugins/
│   └── grafana-app/
│       └── src/
│           ├── api/
│           ├── app/
│           └── features/
│               ├── approvals/
│               ├── canvas/
│               ├── chat/
│               ├── folder-context/
│               ├── knowledge/
│               ├── playbooks/
│               ├── sessions/
│               └── skills/
├── tests/
│   ├── contract/
│   ├── e2e/
│   └── integration/
├── CONTRIBUTING.md
└── README.md
```

## 3. 根级目录说明

| 目录 | 用途 | 内容边界 |
|---|---|---|
| `.github` | GitHub 工程流程配置 | Issue/PR 模板和后续 CI；不放产品实现 |
| `api` | 跨模块、跨语言契约 | HTTP、事件和 MCP 上下文 schema；不放 handler |
| `cmd` | 可执行进程入口 | 配置加载、依赖装配、启动和退出；不放业务规则 |
| `config` | 示例配置 | 仅可提交无密钥的示例；真实 Token、Cookie、Secret 不入库 |
| `data` | Skill、Playbook 和服务种子数据 | 开发环境可直接使用；生产环境应映射持久卷或外部存储 |
| `deploy` | 独立服务的部署材料 | 每个部署单元分别维护；不把四个 MCP Server 合并 |
| `docs` | 用户与团队协作文档 | 不替代正式 Proposal/Design Issue 基线 |
| `internal` | Go 后端的领域和平台实现 | 仓库内部使用，不作为外部 SDK |
| `plugins` | Grafana App Plugin 前端 | UI、前端 SDK 和页面能力，不放 AI Core 业务实现 |
| `tests` | 跨模块验证 | 契约、集成和端到端测试；模块单元测试跟随源码放置 |

## 4. `api`：契约目录

### `api/http`

保存 Chat、Session、Promotion 等 HTTP API 契约。契约按能力拆文件，避免多人同时修改一个总 OpenAPI 文件。

### `api/events`

保存 SSE `AgentEvent`、`AlertEvent` 等跨模块事件格式。事件生产方是对应 schema 的 Owner。

### `api/mcp`

保存 MCP 调用所需的用户、Session、active Folder、追踪信息等公共上下文契约。该目录不定义具体业务工具。

## 5. `cmd`：进程入口

### `cmd/ai-core`

AI Core 主进程入口，装配 Eino Agent、Session、权限、HITL、审计和告警模块。Owner：A。

### `cmd/grafana-app-backend`

Grafana Plugin Backend 入口，负责Grafana身份适配和到 AI Core 的受控转发。Owner：B。

### 四个 `cmd/mcp-*`

| 目录 | 能力 | Owner |
|---|---|---|
| `cmd/mcp-grafana` | Grafana/Prometheus 工具 | C |
| `cmd/mcp-knowledge` | 知识库工具 | D |
| `cmd/mcp-playbook` | Playbook 工具 | E |
| `cmd/mcp-skills` | 对外 Skills 工具 | C |

四个目录对应四个独立 Streamable HTTP MCP Server。每个入口只装配自己的领域 handler，不注册其他 Server 的工具。

## 6. `internal`：后端领域目录

### `internal/agent` — A

- `application`：Agent 应用用例和流程入口。
- `orchestration`：Eino ChatModelAgent、工具选择和上下文编排。
- `transport`：Chat HTTP/SSE 适配。
- `hitl`：Interrupt、CheckPoint、Approve/Reject 和 Resume。

### `internal/audit` — A

审计事件、脱敏和日志写入。其他模块只依赖公开的审计接口，不自行建立另一套审计格式。

### `internal/session` — B

- `canvas`：Session 内的 Canvas 状态模型。
- `store/sqlite/migrations`：Session、Message 和 Canvas 的 SQLite migration。

Session 的 `active_folder_uid` 归此领域维护；Folder Permission 校验仍由平台 Grafana 适配负责。

### `internal/promotion` — B

- `approval`：private → shared 的申请和人工审批状态机。
- `store/sqlite/migrations`：ApprovalRequest 等持久化结构。

`Visibility`、`ApprovalRequest` 由此目录统一定义，Playbook 和 Skill 不复制枚举。

### `internal/grafanaplugin` — B

- `identity`：从 Grafana Plugin Backend 上下文提取并验证用户身份。
- `proxy`：将受控请求转发到 AI Core，避免前端伪造 `user_id`。

### `internal/grafana` — C

Grafana MCP 的业务适配。`mcpserver` 保存 `grafana.*` 工具 schema、handler 与领域服务连接；Grafana Folder Permission 公共客户端不放在这里。

### `internal/knowledge` — D

- `catalog`：ServiceEntry、Runbook、Document 等知识目录能力。
- `retrieval`：按 Folder Scope 执行检索。
- `importer`：后续文档导入边界；MS1可保留为空骨架。
- `mcpserver`：`knowledge.*` MCP 工具适配。
- `store/sqlite/migrations`：知识库持久化 migration。

### `internal/playbook` — E

- `engine`：Eino Graph/Workflow 执行引擎。
- `repository`：Playbook 的读取、保存和可见性过滤接口。
- `mcpserver`：`playbook.*` MCP 工具适配。

### `internal/skills` — C

- `repository`：Skill 的单一逻辑数据源。
- `middleware`：内部 Agent 使用的 Eino Skill Middleware 适配。
- `mcpserver`：外部 AI 工具使用的 Skills MCP 适配。

Middleware 和 MCP Server 必须依赖同一个 repository，不各自维护 Skill 文件。

### `internal/alerting` — E

- `receiver`：P7a Webhook 接收、验签和幂等边界。
- `dispatcher`：P7b Alert → Playbook 映射和异步投递边界。

当前 Proposal 采用 AI Core 内部 channel，因此 MS1 不增加独立 `alert-worker` 进程。

## 7. `internal/platform`：公共技术适配

### `internal/platform/config`

公共配置解析与校验。业务默认值和业务规则仍由所属领域决定。

### `internal/platform/grafana` — D

Grafana API Client、FolderPermissionService、Folder Scope 解析和权限缓存。Knowledge、Playbook、Skills 和 Agent 共用此处的契约。

### `internal/platform/mcp` — C

mcp-go 的 Streamable HTTP 公共装配、认证上下文、错误映射和 read/write 元数据。不注册具体领域工具。

### `internal/platform/storage`

数据库连接、事务和 migration runner 等通用能力。每个领域的表结构仍放在该领域自己的 migration 目录。

## 8. `plugins/grafana-app`：Grafana 前端

### `src/api`

前端业务 SDK 和契约适配层。页面组件不直接散落调用后端接口。

### `src/app`

Grafana App Plugin 注册、路由、全局布局和顶层依赖装配。

### `src/features`

| 目录 | 页面或组件能力 |
|---|---|
| `chat` | 对话输入、SSE消息展示和HITL交互 |
| `canvas` | 图表画布和单图展示 |
| `sessions` | 会话列表、恢复和归档入口 |
| `folder-context` | active Folder选择与上下文显示 |
| `knowledge` | Service、Runbook和Document管理入口 |
| `playbooks` | Playbook列表、详情、运行和编辑入口 |
| `skills` | Skill列表、详情、编辑和运行入口 |
| `approvals` | private → shared申请和审批中心 |

MS1只实现纵向链路需要的最小页面，其余目录用于固定模块边界，不代表MS1要完成全部功能。

## 9. `data`：领域数据根目录

### `data/services`

开发环境的ServiceEntry初始种子。P5正式主存仍通过Repository/DB管理。

### `data/skills`

- `private`：仅Owner可见的Skill。
- `shared`：绑定Grafana Folder的共享Skill。

Eino Skill Middleware和Skills MCP Server共同读取此逻辑数据源。

### `data/playbooks`

- `private`：个人Playbook。
- `shared`：绑定项目Folder或Shared Folder的Playbook。

目录只表达数据组织；最终权限必须通过Owner或Grafana Folder Permission校验，不能依赖文件路径保证安全。

## 10. `deploy`：部署目录

### `deploy/docker`

按进程分别保存容器构建材料。四个MCP Server必须保留四个独立子目录。

### `deploy/k8s`

保存后续Kubernetes部署清单或模板。MS1不要求高可用和复杂编排。

## 11. `tests`：跨模块验证

- `contract`：HTTP、SSE、MCP和领域接口的契约测试。
- `integration`：AI Core、Session、Permission和各MCP Server之间的集成验证。
- `e2e`：从Grafana Plugin发起请求到图表展示的纵向验收。

单个package的单元测试应和源码放在同一目录，不集中堆入`tests`。

## 12. `docs` 与 `.github`

- `docs/PROCESS.md`：GitHub开发流程。
- `docs/DIRECTORY-STRUCTURE.md`：本文，说明目录职责和边界。
- `.github/ISSUE_TEMPLATE`：Issue模板。
- `.github/PULL_REQUEST_TEMPLATE.md`：PR模板。

正式Proposal与Design以GitHub Issue为基线，不在`docs`中维护第二份可漂移的正式规格。

## 13. 禁止出现的目录用法

- 不创建无Owner的`internal/common`、`internal/utils`或`internal/contracts`大杂烩。
- 不把领域逻辑写进`cmd/*/main.go`。
- 不把四个MCP Server合并到一个`cmd`或一个部署单元。
- 不在`config`提交Secret、Token、Cookie或生产连接信息。
- 不在仓库提交构建产物、运行日志、SQLite运行库或编辑器配置。
- 不因目录已经存在，就提前实现不属于当前Milestone的完整功能。
