# 贡献指南

Torchbearing 使用个人 Fork、GitHub Issue 和 Pull Request 协作。开始前请阅读 [GOVERNANCE.md](GOVERNANCE.md)。

## 1. 确认工作类型

### Proposal 草案

产品决策先从正式仓库最新 `main` 创建个人分支，并提交 Draft PR。该 PR 仅用于逐行评审，不合并；内容稳定后，将定稿内容写入处于待决策状态的正式 Proposal Issue，再关闭草案 PR。团队随后在 Issue 中作出并记录接受、拒绝或暂不排期的决定。

正式的功能 Proposal 应表示一个可以独立接受、拒绝或暂不排期，并能通过用户行为验收的产品增量。边界由产品决策决定，不按代码量或 PR 数量划分。跨越多个独立功能的完整体验使用战略方向或 Capability 归类，各功能建立平级 Proposal。

### Research、Roadmap 与前置 Spike

Research、Roadmap 和 Spike 可以在 Proposal 之前直接创建为 Issue：

- Research 提供用户、问题、竞品或资料证据，关闭前记录来源、局限和决策输入；
- Roadmap 维护能力、优先级和候选阶段，不批准其中提到的具体功能；
- Spike 用有时限的实验回答一个主要技术问题，关闭前记录实际方法、结果和结论。

仅当这项研究、规划或实验本身是本轮承诺交付物时才加入 Milestone。

### Design、Spike 与 Task

Design、Spike 和 Task 是按需使用的工程跟踪对象，不是每个 Proposal 的必经层。边界足够小的 Proposal 可以由一个或多个实现 PR 直接关联；存在独立技术取舍、实验风险或多人协作时，再使用：

```text
Design: 描述架构、接口、数据和关键取舍
Spike: 通过实验消除技术不确定性
Task: 定义边界清楚的工程工作
```

创建时在正文用 `#N` 引用上游依据或相关 Issue。GitHub 原生 Sub-issues 仅在确有工作分解、进度汇总或多个负责人时按需使用，不作为 Issue 有效性的前提。横跨多个 Proposal 的工作可以引用多个 Issue，不强制指定唯一父 Issue。Task 使用自动关联的 `task` 类型标签，但不使用 `sub-task` 或 `ms1` 等关系、阶段标签。

Task 只能落实已经确认的产品行为。如果 Task 引入了可独立接受、拒绝、延期或跨 Milestone 排期的用户功能，应暂停工程拆分并创建平级 Proposal。

### Bug

可复现缺陷使用 `Bug:` Issue。涉及凭据、越权或敏感数据的安全问题不得提交普通 Issue，处理方式见 [SECURITY.md](SECURITY.md)。

### 标签、Milestone 与状态

Issue Forms 自动关联与标题前缀一致的稳定类型标签。具有 Triage 及以上权限的成员可以按照治理规则设置、调整其他已有普通标签、Milestone、负责人、Review 请求、Issue 状态和可选的 Sub-issues 关系。`Proposal-Accepted`、`Proposal-Denied`、`Proposal-NoPlan` 只能在团队或指定决策者已经形成并记录结论后设置，不得用权限自行批准 Proposal。

### 完成和关闭 Issue

关闭已完成的 Issue 前，在正文或评论中留下可复核结果。Research 应包含来源、证据、局限和决策输入；Spike 应包含实际环境、方法、结果和结论；Design、Task 等应包含相应决策或验证证据。工作停止或未完成时说明原因，并使用 `Not planned` 关闭理由。

## 2. Fork 与远端

在 GitHub 页面 Fork `1024XEngineer/Torchbearing` 到个人账号，然后克隆个人 Fork：

```bash
git clone https://github.com/<your-name>/Torchbearing.git
cd Torchbearing
git remote add upstream https://github.com/1024XEngineer/Torchbearing.git
git remote -v
```

约定：

- `origin` 指向个人 Fork；
- `upstream` 指向正式仓库；
- 开发分支只存在于个人 Fork；
- 不向正式仓库 `main` 直接推送。

## 3. 从最新 main 创建分支

```bash
git fetch upstream
git switch main
git merge --ff-only upstream/main
git push origin main
```

不要在个人 `main` 上直接开发。为每个 Issue 或 Proposal 草案创建独立分支：

```text
proposal/<topic>
research/<topic>
roadmap/<topic>
design/<issue-number>-<topic>
spike/<issue-number>-<topic>
feat/<issue-number>-<topic>
fix/<issue-number>-<topic>
docs/<issue-number>-<topic>
chore/<issue-number>-<topic>
```

例如：

```bash
git switch -c feat/12-prometheus-query
```

## 4. 提交

提交信息推荐使用：

```text
type(scope): 中文描述
```

常用类型：

| 类型 | 用途 |
| --- | --- |
| `feat` | 新功能 |
| `fix` | 缺陷修复 |
| `refactor` | 不改变行为的重构 |
| `test` | 测试 |
| `docs` | 用户、贡献者或治理文档 |
| `chore` | 依赖、脚本和工程配置 |
| `ci` | 自动化检查配置 |

个人分支可以保留多个便于开发和 Review 的提交。不要为了主线整洁反复强制推送；合并时由 GitHub Squash Merge 统一压缩为一个主线提交。

不得提交 Token、密码、Cookie、私钥、真实凭据、构建产物或本地 `.env`。

## 5. 验证

PR 前执行与改动相关的构建、测试或人工验证，并记录实际命令和结果。

在仓库 CI 尚未建立时，本地验证记录是合并判断的必要证据。CI 建立后，适用检查也必须通过。不得填写未实际执行的命令或伪造结果。

## 6. 创建 Proposal 草案 PR

Proposal 草案 PR 必须：

1. 保持 Draft；
2. 说明评审目标和待决策问题；
3. 明确该 PR 不合并，定稿后写入正式 Issue；
4. 只包含本次讨论需要的文档、原型或 Demo；
5. 通过行内评论保留关键讨论和修改记录。

草案内容稳定后：

1. 创建待决策的正式 Proposal Issue；
2. 从草案复制最终定稿内容；
3. 链接草案 PR；
4. 设置 `proposal` 和适用的规格标签；
5. 关闭草案 PR，不合并；
6. 团队在 Issue 中形成并记录结论后，再设置决策标签和必要的 Milestone。

## 7. 创建可合并 PR

从个人 Fork 向 `1024XEngineer/Torchbearing:main` 创建 PR。正文应：

1. 使用 `Refs #N`、`Part of #N` 或 `Closes #N` 关联主要 Issue；小型 Proposal 可以由 PR 直接关联，不要求先创建 Task；
2. 说明实际改动、依据、影响和明确不做；
3. 列出实际执行的验证命令和结果；
4. 说明 AI 使用和人工复核；
5. 不夹带与主要 Issue 无关的内容；
6. 请求指定维护者进行人工 Review。

只有确实完成整个 Issue 时才使用 `Closes #N`。一个 PR 只完成 Proposal 的一部分时使用 `Part of` 或 `Refs`；存在相关 Design、Spike、Task 等 Issue 时优先关联直接负责的 Issue。

PR 是 Review 的起点，不是合并通知。请及时响应评论，不要自行合并。

## 8. 处理 Review

- 对每条评论回复“已修改”、说明不同意见，或记录为后续 Issue；
- 修改后重新执行受影响的验证；
- 新提交使旧结论失效时主动请求重新 Review；
- 影响项目的关键结论必须保留在 Issue 或 PR 中；
- AI Review 不能代替人工 Review。

## 9. 合并后

可合并 PR 默认由维护者使用 Squash Merge。合并后同步个人 Fork 并删除完成分支：

```bash
git fetch upstream
git switch main
git merge --ff-only upstream/main
git push origin main
git branch -d feat/12-prometheus-query
git push origin --delete feat/12-prometheus-query
```

## 10. 文档与 AI

- 产品和架构工程文档定稿在 Issue 中，不作为普通文档合入仓库；
- README、用户文档、部署说明和治理文件随仓库维护；
- 设计进入编码后视为基线，变化应新建 Issue；
- AI 生成内容必须由作者理解、验证并承担维护责任；
- 不向外部 AI 服务提供凭据、内部地址或未脱敏数据。
