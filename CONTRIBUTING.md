# 贡献指南

## 1. 开始之前

所有实现都应从一个边界清楚的 Issue 开始。Issue 至少说明：

- 为什么要做；
- 本次做什么；
- 明确不做什么；
- 怎样验证完成。

较大的产品变化先提交 Proposal，待标记为 `Proposal-Accepted` 后再进入 Design 和开发。

## 2. 分支与提交

- 从最新 `main` 创建分支；
- 推荐命名：`agent/<description>`、`feat/<description>`、`fix/<description>`；
- 一个分支只处理一个 Issue；
- 提交信息简短描述实际增量；
- 不提交密钥、Token、构建产物或本地配置。

## 3. Pull Request

PR 必须：

1. 使用 `Closes #N` 或明确关联 Issue；
2. 说明改动、原因和影响；
3. 列出实际执行的验证；
4. 不夹带无关修改；
5. 保持规模可评审；
6. 经人工 Review 后合并。

AI 参与的代码必须有人能够解释、测试和维护。

## 4. 文档

- 产品/架构草案 PR 只用于讨论，不合并；
- 定稿内容落回对应 Issue；
- 用户文档和贡献文档通过普通 PR 合并；
- 功能已完成但缺少用户文档时使用 `Need-Document`；
- 文档完成后切换为 `Documented`。

## 5. Milestone 完成

每个 Milestone 结束前检查：

- 交付 Issue 已完成或明确移出；
- 主干可构建和演示；
- 关键测试通过；
- 架构变化有记录；
- Release 描述列出功能、限制和已知问题。
