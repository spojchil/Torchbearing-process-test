# Torchbearing Process Test

这是 Torchbearing 项目的 GitHub 工程流程测试仓库，用来验证：

```text
Milestone → Proposal → Design → Task → PR → Review → Merge → Release
```

当前产品方向是 Grafana 中的自然语言 Prometheus 指标分析工作台。MS1 以“成功生成一张真实 Prometheus 临时图”为闭环。

## 工作方式

- Milestone 定义每个两周迭代的目标和交付清单；
- Proposal Issue 记录产品取舍和验收边界；
- Design Issue 固定模块、接口、数据结构和风险；
- Task/sub-task Issue 对应一个或少量可评审 PR；
- 每个 Milestone 结束发布 Tag 和 Release。

详细规则见 [GitHub 开发流程](docs/PROCESS.md) 和 [贡献指南](CONTRIBUTING.md)。

## 当前 Milestone

- MS1：Prometheus 单图闭环
- MS2：可持续指标分析
- MS3：经验沉淀与复用
- MS4：冻结、质量与发布

具体范围和完成度以仓库的 [Milestones](https://github.com/spojchil/Torchbearing-process-test/milestones) 为准。

## 文档归属

- 产品和架构工程文档以 Issue 为正式基线；
- 草案可以使用不合并的 Draft PR 逐行讨论；
- README、用户手册和贡献指南随代码进入仓库；
- 已接受决策发生变化时新增 Issue，不覆盖原始决策记录。
Torchbearing 产品、架构与 GitHub Milestone 流程测试仓库
