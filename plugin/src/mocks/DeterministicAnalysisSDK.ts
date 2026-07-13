import type { AnalysisRequest, AnalysisResponse } from '../contracts/analysis';
import { AnalysisSDKError, ErrorCodes, type ErrorResponse } from '../contracts/errors';
import type { AnalysisSDK } from '../sdk/AnalysisSDK';
import {
  getFixtureError,
  getFixtureResponse,
  type MockScenario,
} from './fixtureRegistry';

/** DeterministicAnalysisSDK 使用固定 fixture 实现前端 AnalysisSDK，不发起任何网络请求。 */
export class DeterministicAnalysisSDK implements AnalysisSDK {
  public constructor(private readonly scenario: MockScenario) {}

  public async analyze(request: AnalysisRequest): Promise<AnalysisResponse> {
    const requestError = validateRequest(request);
    if (requestError !== undefined) {
      throw new AnalysisSDKError(requestError);
    }

    switch (this.scenario) {
      case 'success':
      case 'empty':
        return getFixtureResponse(this.scenario);
      case 'agent-failure':
      case 'metrics-failure':
      case 'invalid-scope':
        throw new AnalysisSDKError(getFixtureError(this.scenario));
      default:
        throw unsupportedScenario(this.scenario);
    }
  }
}

function validateRequest(request: AnalysisRequest): ErrorResponse | undefined {
  if (request.text.trim() === '') {
    return {
      code: ErrorCodes.invalidArgument,
      message: 'analysis text is required',
      retryable: false,
      requestId: 'mock-analysis-boundary',
    };
  }
  if (
    request.scope.datasourceUid.trim() === '' ||
    request.scope.timeRange.from.trim() === '' ||
    request.scope.timeRange.to.trim() === ''
  ) {
    return {
      code: ErrorCodes.invalidScope,
      message: 'complete analysis scope is required',
      retryable: false,
      requestId: 'mock-analysis-boundary',
    };
  }
  return undefined;
}

function unsupportedScenario(scenario: never): AnalysisSDKError {
  return new AnalysisSDKError({
    code: ErrorCodes.internal,
    message: `unsupported deterministic frontend scenario: ${String(scenario)}`,
    retryable: false,
    requestId: 'mock-analysis-unsupported',
  });
}
