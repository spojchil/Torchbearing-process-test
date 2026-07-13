import { ErrorCodes } from '../../contracts/errors';
import { DeterministicAnalysisSDK } from '../../mocks/DeterministicAnalysisSDK';
import { createDefaultScopeDraft, toAnalysisScope } from '../scope/scopeModel';
import { AnalysisWorkbench } from './AnalysisWorkbench';

// @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
describe('AnalysisWorkbench', () => {
  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('把成功响应转换为 success 状态', async () => {
    const workbench = new AnalysisWorkbench(new DeterministicAnalysisSDK('success'));
    const state = await workbench.submit('查看 checkout 请求速率', defaultScope());

    assertEqual(state.status, 'success');
    assertEqual(state.charts.length, 1);
    assertEqual(state.error, undefined);
  });

  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('把空图表响应转换为 empty 状态', async () => {
    const workbench = new AnalysisWorkbench(new DeterministicAnalysisSDK('empty'));
    const state = await workbench.submit('查看未知服务', defaultScope());

    assertEqual(state.status, 'empty');
    assertEqual(state.charts.length, 0);
  });

  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('把 SDK typed error 转换为 error 状态', async () => {
    const workbench = new AnalysisWorkbench(new DeterministicAnalysisSDK('metrics-failure'));
    const state = await workbench.submit('查看 checkout 请求速率', defaultScope());

    assertEqual(state.status, 'error');
    assertEqual(state.error?.code, ErrorCodes.metricsUnavailable);
    assertEqual(state.error?.retryable, true);
  });

  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('返回状态副本并支持确定性重置', async () => {
    const workbench = new AnalysisWorkbench(new DeterministicAnalysisSDK('success'));
    const state = await workbench.submit('查看 checkout 请求速率', defaultScope());
    state.charts.length = 0;

    assertEqual(workbench.getState().charts.length, 1);
    assertEqual(workbench.reset().status, 'idle');
  });
});

function defaultScope() {
  return toAnalysisScope(createDefaultScopeDraft());
}

function assertEqual(actual: unknown, expected: unknown): void {
  if (actual !== expected) {
    throw new Error(`actual ${String(actual)} does not equal expected ${String(expected)}`);
  }
}
