# Torchbearing

Torchbearing 是面向 Grafana 可观测性分析场景的 AI 辅助项目。

项目当前处于初始化阶段。仓库治理基线、代码骨架和自动化检查分别通过独立 Pull Request 评审；在正式内容合入前，以 GitHub Issue、Pull Request 和 Milestone 中的记录为准。

## 协作与治理

项目采用 Issue 驱动、个人 Fork 开发、Pull Request 评审和 Squash Merge 的协作方式。

- [项目治理规则](GOVERNANCE.md)：对象职责、Proposal 流程、Milestone、Review 和合并规则；
- [贡献指南](CONTRIBUTING.md)：Fork、分支、提交、验证和 PR 的具体操作；
- [安全策略](SECURITY.md)：安全问题报告、凭据和数据处理要求；
- [团队行为准则](CODE_OF_CONDUCT.md)：协作、分歧处理和行为边界。

核心过程如下：

```text
Research / Roadmap / 前置 Spike（按需）
→ Proposal 草案 Draft PR（逐行评审，不合并）
→ 待决策的正式 Proposal Issue
→ 团队记录接受 / 拒绝 / 暂不排期
  ├─ 边界足够小：代码 PR 直接关联已接受的 Proposal
  └─ 工作较复杂：按需建立 Design / Spike / Task → 代码 PR
→ 验证与 Review
→ Squash Merge
→ Milestone Release
```
