import type { AnalysisRequest } from '../contracts/analysis';
import { AnalysisSDKError, ErrorCodes, type ErrorCode } from '../contracts/errors';
import { DeterministicAnalysisSDK } from './DeterministicAnalysisSDK';

// @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
describe('DeterministicAnalysisSDK', () => {
  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('稳定返回成功 fixture 的独立副本', async () => {
    const sdk = new DeterministicAnalysisSDK('success');
    const first = await sdk.analyze(validRequest());
    const second = await sdk.analyze(validRequest());

    assertEqual(first.requestId, 'mock-analysis-001');
    assertEqual(first.charts.length, 1);
    assert(first !== second, '每次调用必须返回独立响应对象');
    assert(first.charts !== second.charts, '每次调用必须返回独立图表数组');
  });

  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('稳定返回空数据 fixture', async () => {
    const response = await new DeterministicAnalysisSDK('empty').analyze(validRequest());

    assertEqual(response.requestId, 'mock-analysis-002');
    assertEqual(response.charts.length, 0);
  });

  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('稳定复现 Agent 和 Metrics 失败', async () => {
    await assertSDKError('agent-failure', ErrorCodes.agentUnavailable);
    await assertSDKError('metrics-failure', ErrorCodes.metricsUnavailable);
  });

  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('稳定复现无效范围边界', async () => {
    await assertSDKError('invalid-scope', ErrorCodes.invalidScope);
  });

  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('拒绝空文本边界输入', async () => {
    const request = validRequest();
    request.text = '  ';

    try {
      await new DeterministicAnalysisSDK('success').analyze(request);
      throw new Error('expected SDK to reject empty text');
    } catch (error: unknown) {
      assert(error instanceof AnalysisSDKError, '错误必须是 AnalysisSDKError');
      assertEqual(error.details.code, ErrorCodes.invalidArgument);
    }
  });
});

async function assertSDKError(
  scenario: 'agent-failure' | 'metrics-failure' | 'invalid-scope',
  code: ErrorCode
): Promise<void> {
  try {
    await new DeterministicAnalysisSDK(scenario).analyze(validRequest());
    throw new Error(`expected ${scenario} to fail`);
  } catch (error: unknown) {
    assert(error instanceof AnalysisSDKError, '错误必须是 AnalysisSDKError');
    assertEqual(error.details.code, code);
  }
}

function validRequest(): AnalysisRequest {
  return {
    text: '查看 checkout 服务过去 30 分钟的请求速率',
    scope: {
      datasourceUid: 'prometheus-mock',
      timeRange: { from: 'now-30m', to: 'now' },
    },
  };
}

function assert(condition: boolean, message: string): asserts condition {
  if (!condition) {
    throw new Error(message);
  }
}

function assertEqual(actual: unknown, expected: unknown): void {
  if (actual !== expected) {
    throw new Error(`actual ${String(actual)} does not equal expected ${String(expected)}`);
  }
}
