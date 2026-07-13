// @ts-ignore -- 依赖已在共享 manifest 中声明，但当前工作区尚未安装类型包。
import { AppPlugin } from '@grafana/data';
// @ts-ignore -- 依赖已在共享 manifest 中声明，但当前工作区尚未安装类型包。
import { createElement, useMemo, useState } from 'react';
import type { ChartSpec } from './contracts/analysis';
import { AnalysisWorkbench, type AnalysisWorkbenchState } from './features/analysis/AnalysisWorkbench';
import { createDefaultScopeDraft, describeScope, toAnalysisScope } from './features/scope/scopeModel';
import { DeterministicAnalysisSDK } from './mocks/DeterministicAnalysisSDK';
import type { MockScenario } from './mocks/fixtureRegistry';

interface TextInputEvent {
  currentTarget: { value: string };
}

/** createMS1Workbench 只通过 AnalysisSDK 组装 D 的工作台，不访问后端内部实现。 */
export function createMS1Workbench(scenario: MockScenario = 'success'): AnalysisWorkbench {
  return new AnalysisWorkbench(new DeterministicAnalysisSDK(scenario));
}

/** AnalysisRootPage 是 MS1 的最小 Grafana 页面，只展示确定性 mock 主流程。 */
function AnalysisRootPage() {
  const workbench = useMemo(() => createMS1Workbench('success'), []) as AnalysisWorkbench;
  const scope = useMemo(() => toAnalysisScope(createDefaultScopeDraft()), []);
  const [text, setText] = useState('查看 checkout 服务过去 30 分钟的请求速率') as [
    string,
    (next: string) => void,
  ];
  const [state, setState] = useState(workbench.getState()) as [
    AnalysisWorkbenchState,
    (next: AnalysisWorkbenchState) => void,
  ];

  const runAnalysis = async () => {
    const pending = workbench.submit(text, scope);
    setState(workbench.getState());
    setState(await pending);
  };

  const chartNodes = state.charts.map((chart: ChartSpec) =>
    createElement(
      'article',
      { key: chart.id },
      createElement('h3', null, chart.title),
      createElement('p', null, `${chart.type} · ${chart.datasourceUid}`),
      createElement('code', null, chart.promql)
    )
  );

  return createElement(
    'main',
    { 'aria-label': 'Torchbearing MS1 分析工作台' },
    createElement('h1', null, 'Torchbearing 指标分析工作台'),
    createElement('p', null, describeScope(scope)),
    createElement('label', null, '分析问题'),
    createElement('input', {
      'aria-label': '分析问题',
      value: text,
      onChange: (event: TextInputEvent) => setText(event.currentTarget.value),
    }),
    createElement(
      'button',
      { disabled: state.status === 'loading', onClick: () => void runAnalysis() },
      state.status === 'loading' ? '分析中…' : '开始分析'
    ),
    createElement('p', { role: state.status === 'error' ? 'alert' : 'status' }, state.message),
    state.error === undefined ? null : createElement('code', null, state.error.code),
    createElement('section', { 'aria-label': '分析图表' }, ...chartNodes)
  );
}

// Grafana 在加载插件时读取该导出并挂载根页面。
export const plugin = new AppPlugin().setRootPage(AnalysisRootPage);
