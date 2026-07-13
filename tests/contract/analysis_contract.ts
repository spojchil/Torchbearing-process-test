// 本文件通过 TypeScript 编译期检查，确保前端请求、响应、错误和 SDK 接口保持一致。
import type {
  AnalysisRequest,
  AnalysisResponse,
} from '../../plugin/src/contracts/analysis';
import { AnalysisSDKError, ErrorCodes } from '../../plugin/src/contracts/errors';
import type { AnalysisSDK } from '../../plugin/src/sdk/AnalysisSDK';

const request: AnalysisRequest = {
  text: '查看 checkout 服务过去 30 分钟的请求速率',
  scope: {
    datasourceUid: 'prometheus-mock',
    timeRange: { from: 'now-30m', to: 'now' },
  },
};

const response: AnalysisResponse = {
  requestId: 'mock-analysis-001',
  message: '已生成 checkout 服务请求速率图表。',
  charts: [],
  mock: true,
};

const deterministicSDK: AnalysisSDK = {
  analyze: async (input) => {
    if (input.text === '') {
      throw new AnalysisSDKError({
        code: ErrorCodes.invalidArgument,
        message: 'text is required',
        retryable: false,
        requestId: 'mock-analysis-boundary',
      });
    }
    return response;
  },
};

void deterministicSDK.analyze(request);
