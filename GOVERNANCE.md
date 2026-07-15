# 项目治理规则

本文定义 Torchbearing 使用 GitHub 进行产品决策、任务管理、代码评审和发布的规则。具体 Git 操作见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 1. 基本原则

- 产品和架构先于编码，重要结论不得只保留在聊天记录中；
- Milestone 管本轮 Issue，Issue 管决策与任务，Pull Request 管评审和代码；
- 开发者通过个人 Fork 工作，不直接向正式仓库 `main` 推送；
- AI 可以辅助设计、编码和 Review，但最终判断和责任由人承担；
- 仓库主线默认使用 Squash Merge，使一个可追溯的 PR 对应一个主线提交。

## 2. 角色与权限

项目遵循最小权限原则。

| 角色 | 仓库权限 | 主要职责 |
| --- | --- | --- |
| 开发者 | Triage | 通过个人 Fork 开发，创建和整理 Issue、创建 PR、响应 Review |
| 仓库协作者 | Write | 日常 Review、合并准备、过程维护 |
| 仓库创建者 | Admin | 仓库设置、权限、最终合并与发布管理 |

本文所称“维护者”包括具有 Write 或 Admin 权限的仓库协作者和仓库创建者。具有 Triage 及以上权限的成员可以按照本规则设置已有的普通类型标签、Milestone、负责人、Review 请求、Issue 状态和可选的原生 Sub-issues 关系。

`Proposal-Accepted`、`Proposal-Denied` 和 `Proposal-NoPlan` 是团队决策的记录，不是普通整理标签。它们只能在团队或指定决策者已经形成并记录结论后设置；Triage 权限不代表可以自行批准自己的 Proposal。

开发分支只存在于个人 Fork。任何人不得绕过 PR 直接修改正式仓库 `main`。

## 3. GitHub 对象职责

### 3.1 Milestone

每个 MS 建立一个 Milestone，只收纳本轮承诺完成的 Issue，不收纳 PR。Issue 创建者可以在正文中说明建议阶段，具有 Triage 及以上权限的成员根据已经确认的计划设置或调整 Milestone。

Milestone 描述应包含本轮目标、验收条件、主要功能和明确不做。长期 Capability、尚未排期的 Proposal 和普通待办不进入当前 Milestone。Research、Roadmap 或 Spike 只有在该项工作本身是本轮承诺交付物时才进入 Milestone；进入 Milestone 不代表其中提到的产品功能已经被接受。

不使用 `ms1`、`ms2` 等阶段标签；阶段归属只由 Milestone 表达。

### 3.2 Capability

Capability 表示跨 Milestone 的长期产品能力，用 `Capability:` 标题前缀识别。它用于维护产品能力地图，不直接代表本轮交付，也不默认进入 Milestone。

### 3.3 Research、Roadmap 与 Spike

Research、Roadmap 和 Spike 可以出现在正式 Proposal 之前，为产品决策提供输入：

- `Research:` 记录用户问题、竞品、资料或事实证据，并说明来源、局限和对决策的影响；
- `Roadmap:` 作为可持续更新的索引，维护能力、优先级和候选阶段的地图，不代替具体功能的产品决策；
- `Spike:` 用有时限的实验回答一个主要技术问题，关闭前必须补充真实环境、方法、结果和结论。

它们可以直接创建为 Issue，也可以关联草案 PR、Proposal 或 Design。Research 和 Roadmap 中出现某项功能，不表示该功能已经获批；需要独立接受或否决的功能仍应建立 Proposal。

### 3.4 Issue

Issue 是正式决策和开发任务的基本单元。标题使用以下前缀：

```text
Proposal:
Research:
Roadmap:
Design:
Spike:
Task:
Bug:
Documentation:
Capability:
```

存在上游依据或相关工作的 Issue，正文必须用 `#N` 明确引用对应 Issue 或 PR。GitHub 原生 Sub-issues 是可选的进度组织工具，只在确有工作分解、进度汇总或多个负责人时使用；跨多个 Proposal 的横切工作使用正文引用，不强制指定唯一父 Issue。未设置原生父子关系不影响追溯有效性。

`Task:` 使用 `task` 类型标签，但不使用 `sub-task` 标签。使用原生 Sub-issues 时，不在正文中重复维护完整的子任务清单。

一个 Issue 应当背景清楚、目标明确、范围可控、验收标准可验证。需要进入仓库的实现或文档工作通常对应一个或少量 PR；Research、Roadmap 和纯决策 Issue 可以不产生 PR。

Issue 完成后应补充可复核的关闭证据：Research 记录来源、证据、局限和决策输入；Spike 记录实际实验环境、方法、结果和结论；Design、Task 等记录相应决策或验证结果。工作停止或未完成时记录原因并使用 GitHub 的 `Not planned` 关闭理由。快速关闭本身不是问题，缺少完成依据才是问题。

### 3.5 Pull Request

PR 分为两类：

1. **设计草案 PR**：用于 Proposal、产品设计或原型的逐行讨论，保持 Draft，最终关闭且不合并；
2. **可合并 PR**：承载代码、治理、用户文档或其他应进入仓库的变更。

PR 不设置 Milestone。产品决策状态、规格粒度和阶段归属由 Issue 管理。

## 4. Proposal 流程

### 4.1 草案讨论

产品决策先建立 Proposal 草案 Draft PR。草案 PR 正文必须明确：

- 该 PR 仅用于逐行评审，不执行合并；
- 定稿后以正式 Proposal Issue 为基线；
- 后续变化通过新的 Proposal 记录，不覆盖旧决策。

维护者可以为草案 PR 设置 `proposal`、`design`、`architecture`、`research`、`prototype` 等标签辅助检索，但不得用草案 PR 代替最终 Issue。

### 4.2 定稿为 Issue

草案内容稳定到足以整体评审后，将定稿内容写入正式 Proposal Issue，并在 Issue 中链接草案 PR。随后关闭草案 PR，不合并。此时 Proposal 默认处于待决策状态；创建正式 Issue 不表示已经接受。

Proposal Issue 使用：

- 类型：`proposal`；
- 决策：待团队形成并记录结论后，设置 `Proposal-Accepted`、`Proposal-Denied` 或 `Proposal-NoPlan`；
- 规格：根据影响范围使用 `MiniSpec` 或 `FullSpec`；被拒绝或暂不排期时可以移除规格标签。

被拒绝或暂不排期的 Proposal 记录理由后关闭，不进入当前 Milestone。已接受且承诺本轮交付的 Proposal 进入对应 Milestone。

战略 Proposal 的交付物是方向决策，正式基线记录完成后即可关闭。功能 Proposal 通常保持开放，直到整体产品验收标准满足。

### 4.3 Proposal 粒度

功能 Proposal 是可以被团队独立接受、拒绝或暂不排期的最小产品增量。Proposal 边界由产品决策决定，不由代码量、模块数量、PR 数量或开发人数决定。

判断一项内容是否应成为独立 Proposal，主要看它是否：

- 形成独立、可感知的用户行为或价值；
- 可以在不否定其他能力的情况下被单独接受、拒绝或延期；
- 可能进入不同 Milestone；
- 需要自己的范围、明确不做和产品验收标准。

如果一个“完整体验”包含多个满足上述条件的功能，它应作为战略方向或 Capability 进行归类，各功能分别建立平级 Proposal。接受战略方向不等于自动接受其中所有功能 Proposal。

Proposal 下面不再创建 Proposal。新的产品决策建立平级 Proposal，并在正文中引用相关战略、Capability 或既有 Proposal。

### 4.4 实现路径

Design、Spike 和 Task 不是 Proposal 到代码 PR 的必经层。已接受的功能 Proposal 根据实际复杂度选择路径：

```text
Proposal
├── 边界足够小：一个或多个实现 PR 直接关联 Proposal
└── 工作较复杂：按需建立以下相关 Issue，再由 PR 关联相应 Issue
    ├── Design
    ├── Spike
    ├── Task
    ├── Bug
    └── Documentation
```

以下情况才需要单独创建 Design、Spike 或 Task：

- 存在需要独立评审的架构、接口或数据取舍；
- 存在需要通过实验消除的关键不确定性；
- 需要多人协作、多个 PR、独立负责人或单独进度跟踪。

Task 只负责落实已经确认的产品行为，不得借工程拆分引入新的用户功能或扩大 Proposal 边界。如果 Task 中出现可以独立接受、拒绝或排期的产品行为，应停止实现并建立平级 Proposal。

相关 Issue 可以仅通过正文中的 `#N` 建立可追溯关系；需要汇总进度时再使用原生 Sub-issues。Design、Spike 和 Task 不因未设置原生父 Issue 而失效。

## 5. 文档归属

- 产品 Proposal、架构设计和实现澄清属于工程文档，定稿载体是 Issue；
- 工程文档进入编码后视为只读基线，变化应创建新 Issue 记录原因和影响；
- Draft PR 只保留讨论和演进历史，不将草案文件合入仓库；
- README、用户使用说明、部署说明、贡献指南和治理文件进入仓库并随项目维护；
- 架构空骨架是代码，经 Review 后应合入仓库，不属于“不合并的设计草案”。

## 6. PR 关联规则

每个可合并 PR 必须关联一个主要 Issue：

- `Refs #N`：提供关联信息，但不表示完成整个 Issue；
- `Part of #N`：实现多 PR 工作中的一部分；
- `Closes #N`：该 PR 合并后确实完成整个 Issue。

边界足够小的 Proposal 可以作为实现 PR 的主要 Issue，不要求先创建 Task。复杂工作存在相关 Issue 时，PR 优先关联直接负责的 Design、Spike、Task、Bug 或 Documentation，并通过正文引用或可选的 Sub-issues 关系追溯到 Proposal。

不得为了自动关闭而对未完成的父 Proposal 使用 `Closes`。一个 PR 原则上只处理一个主要 Issue，不夹带无关改动。

设计草案 PR 出现在正式 Proposal Issue 之前，因此不要求关联一个已有的 Proposal Issue；定稿 Issue 必须反向链接草案 PR。

## 7. Review 与合并

每个可合并 PR 都应有 Author、Reviewer 和 Merger。

- Author 对范围、实现、验证和说明负责；
- Reviewer 原则上由仓库协作者担任，检查方向、范围、实现、验证和风险；
- Merger 原则上由仓库创建者担任，在确认 Review 结论和合并条件后执行合并。

Reviewer 和 Merger 原则上由不同的人承担。作者不得批准自己的 PR；AI Review 不计作人工 Review。

可合并 PR 必须满足：

1. 已关联主要 Issue，范围与 Issue 一致；
2. 作者能够说明做了什么、为什么这样做以及明确不做什么；
3. 与改动相关的构建、测试或人工验证已经完成；
4. 至少完成一次人工 Review，阻塞问题已经处理；
5. 作者能够解释和维护包括 AI 生成内容在内的全部变更；
6. 不包含凭据、敏感数据、调试残留或本地配置。

可合并 PR 默认使用 GitHub Squash Merge。个人分支允许保留多个便于开发和 Review 的提交，不要求开发者为主线整洁反复改写分支历史。

## 8. 自动化检查

在仓库 CI 尚未建立时，PR 作者必须列出实际执行的本地构建、测试和人工验证命令及结果，由 Reviewer 判断证据是否充分。

CI 建立后，所有适用检查必须通过才能合并。Fork PR 工作流不得获得仓库 Secrets 或写 Token。CI 配置及其变更通过独立 PR 评审。

## 9. Milestone 与 Release

Milestone 结束时：

1. 核对承诺 Issue 的状态和验收结果；
2. 记录未完成项、变更原因和后续去向；
3. 形成可打 Tag 的版本；
4. 创建 GitHub Release，列出用户可感知功能、修复、限制和已知问题。

PR 的合并数量不是 Milestone 完成度；完成度以 Issue 为准。

## 10. 标签使用

标签用于表达需要跨 Issue 检索的类型、决策和状态，不重复表达标题、Sub-issues 或 Milestone 已有的信息。

Issue Forms 自动添加与标题前缀一致的稳定类型标签。具有 Triage 及以上权限的成员可以根据正文和当前计划设置、调整其他普通标签、Milestone 与 Issue 状态。Proposal 决策标签只能在人工决策已经形成并记录后设置，作者不得用标签自行批准 Proposal。

- Proposal：`proposal`；
- 长期能力：`capability`；
- 前置证据与规划：`research`、`roadmap`、`spike`；
- 决策：`Proposal-Accepted`、`Proposal-Denied`、`Proposal-NoPlan`；
- 规格：`MiniSpec`、`FullSpec`；
- 工程与设计：`design`、`architecture`、`task`、`prototype`；
- 缺陷：`bug`；
- 文档：`documentation`、`Need-Document`、`Documented`；
- 通用处置：`duplicate`、`invalid`。

`documentation` 只用于需要进入仓库的用户或贡献者文档，不用于产品和架构草案。

## 11. AI 使用

AI 可以辅助调研、设计、编码、测试和 Review，但不能代替人的判断和责任。

- 关键变更应说明 AI 的生成意图、关键输入和人工复核内容；
- 作者必须理解控制流、数据结构、错误处理和安全边界；
- 无法解释、验证或维护的内容不得合并；
- 不得向外部 AI 服务提交 Token、凭据或未脱敏的内部数据。

## 12. 规则变更

治理规则变更必须通过 Issue 说明动机和影响，并通过可合并 PR Review 后进入仓库。不得只通过聊天记录改变正式规则。
