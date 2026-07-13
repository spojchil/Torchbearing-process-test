import agentFailureFixture from '../../../contracts/fixtures/agent-failure.json';
import emptyFixture from '../../../contracts/fixtures/empty.json';
import invalidScopeFixture from '../../../contracts/fixtures/invalid-scope.json';
import metricsFailureFixture from '../../../contracts/fixtures/metrics-failure.json';
import successFixture from '../../../contracts/fixtures/success.json';
import type { AnalysisResponse, ChartSpec, ChartType } from '../contracts/analysis';
import { ErrorCodes, type ErrorCode, type ErrorResponse } from '../contracts/errors';

/** MockScenario 表示前端 MS1 mock 支持的固定场景。 */
export type MockScenario =
  | 'success'
  | 'empty'
  | 'agent-failure'
  | 'metrics-failure'
  | 'invalid-scope';

export type SuccessScenario = Extract<MockScenario, 'success' | 'empty'>;
export type FailureScenario = Exclude<MockScenario, SuccessScenario>;

interface FixtureChart {
  id: string;
  title: string;
  type: string;
  datasourceUid: string;
  promql: string;
  timeRange: { from: string; to: string };
}

interface ResponseFixture {
  response: {
    requestId: string;
    message: string;
    charts: FixtureChart[];
    mock: boolean;
  };
}

interface ErrorFixture {
  error: {
    code: string;
    message: string;
    retryable: boolean;
    requestId: string;
  };
}

const successResponse = readResponseFixture(successFixture);
const emptyResponse = readResponseFixture(emptyFixture);
const agentFailure = readErrorFixture(agentFailureFixture);
const metricsFailure = readErrorFixture(metricsFailureFixture);
const invalidScope = readErrorFixture(invalidScopeFixture);

/** getFixtureResponse 返回成功或空数据场景的独立副本，避免测试间共享可变状态。 */
export function getFixtureResponse(scenario: SuccessScenario): AnalysisResponse {
  switch (scenario) {
    case 'success':
      return cloneResponse(successResponse);
    case 'empty':
      return cloneResponse(emptyResponse);
  }
}

/** getFixtureError 返回确定性失败场景的独立错误副本。 */
export function getFixtureError(scenario: FailureScenario): ErrorResponse {
  switch (scenario) {
    case 'agent-failure':
      return { ...agentFailure };
    case 'metrics-failure':
      return { ...metricsFailure };
    case 'invalid-scope':
      return { ...invalidScope };
  }
}

function readResponseFixture(fixture: ResponseFixture): AnalysisResponse {
  if (fixture.response.mock !== true) {
    throw new Error('MS1 response fixture must be marked as mock');
  }

  return {
    requestId: fixture.response.requestId,
    message: fixture.response.message,
    charts: fixture.response.charts.map(readChart),
    mock: true,
  };
}

function readChart(chart: FixtureChart): ChartSpec {
  return {
    id: chart.id,
    title: chart.title,
    type: readChartType(chart.type),
    datasourceUid: chart.datasourceUid,
    promql: chart.promql,
    timeRange: { ...chart.timeRange },
  };
}

function readChartType(value: string): ChartType {
  switch (value) {
    case 'timeseries':
    case 'stat':
    case 'table':
      return value;
    default:
      throw new Error(`unsupported chart type in fixture: ${value}`);
  }
}

function readErrorFixture(fixture: ErrorFixture): ErrorResponse {
  return {
    code: readErrorCode(fixture.error.code),
    message: fixture.error.message,
    retryable: fixture.error.retryable,
    requestId: fixture.error.requestId,
  };
}

function readErrorCode(value: string): ErrorCode {
  switch (value) {
    case ErrorCodes.invalidArgument:
    case ErrorCodes.invalidScope:
    case ErrorCodes.agentUnavailable:
    case ErrorCodes.metricsUnavailable:
    case ErrorCodes.noData:
    case ErrorCodes.internal:
      return value;
    default:
      throw new Error(`unsupported error code in fixture: ${value}`);
  }
}

function cloneResponse(response: AnalysisResponse): AnalysisResponse {
  return {
    ...response,
    charts: response.charts.map((chart) => ({
      ...chart,
      timeRange: { ...chart.timeRange },
    })),
  };
}
