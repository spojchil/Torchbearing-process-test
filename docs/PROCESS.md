# GitHub 开发流程

## 1. 对象职责

| 对象 | 职责 |
|---|---|
| Milestone | 定义一个迭代的目标、截止日期和能力清单 |
| Proposal Issue | 记录产品决策和边界 |
| Design Issue | 固定架构、接口、数据结构和风险 |
| Task Issue | 承载可执行工程增量 |
| Pull Request | 承载代码、验证与 Review |
| Release | 固化该轮可交付版本 |

## 2. 标准链路

```text
Proposal
  → Proposal-Accepted / Proposal-Denied / Proposal-NoPlan
  → Design
  → Task / sub-task
  → Pull Request
  → Review
  → Merge
  → Tag / Release
```

## 3. Proposal 状态

| 标签 | 含义 |
|---|---|
| `proposal` | 仍在讨论 |
| `Proposal-Accepted` | 已接受，可以进入实现 |
| `Proposal-Denied` | 已拒绝，保留理由后关闭 |
| `Proposal-NoPlan` | 方向成立，但当前版本不排期 |
| `FullSpec` | 影响面较大的完整规格 |
| `MiniSpec` | 小范围精简规格 |

`Proposal-NoPlan` 不应归属当前 Milestone。Proposal 一旦进入开发即作为当时基线，不回头覆盖；变化通过新 Issue 描述。

## 4. Task 和 PR

Task 应尽量对应一个 PR。确需多个 PR 时，应在 Issue 中列出阶段和依赖。

PR 合并前检查：

- [ ] 关联 Issue；
- [ ] 改动范围与 Issue 一致；
- [ ] 构建、测试或人工验证通过；
- [ ] 没有无关文件；
- [ ] Review 意见已处理；
- [ ] 用户文档状态正确。

## 5. 草案文档

产品、架构和原型可以使用 Draft PR 获得逐行评论，但这类 PR 不合并。定稿内容复制回 Issue，并在 Issue 中保留草案 PR 链接。

## 6. Milestone 和 Release

每个 Milestone 说明应直接列出能力，而不是堆放过程描述。结束时：

1. 完成或移出未交付 Issue；
2. 记录架构演进；
3. 发布 Tag/Release；
4. 在 Release 中列出功能、限制、已知问题和验证结果；
5. 明确下一 Milestone 的输入。

## 7. 参考流程

- [XGo Milestones](https://github.com/goplus/xgo/milestones)
- [Accepted FullSpec 示例](https://github.com/goplus/xgo/issues/2802)
- [Accepted MiniSpec 示例](https://github.com/goplus/xgo/issues/2751)
- [Denied Proposal 示例](https://github.com/goplus/xgo/issues/2667)
