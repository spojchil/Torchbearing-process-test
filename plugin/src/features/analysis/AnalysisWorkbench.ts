import type { AnalysisScope, ChartSpec } from '../../contracts/analysis';
import { AnalysisSDKError, ErrorCodes, type ErrorResponse } from '../../contracts/errors';
import type { AnalysisSDK } from '../../sdk/AnalysisSDK';

export type WorkbenchStatus = 'idle' | 'loading' | 'success' | 'empty' | 'error';

/** AnalysisWorkbenchState 是页面组件消费的稳定展示状态。 */
export interface AnalysisWorkbenchState {
  status: WorkbenchStatus;
  message: string;
  charts: ChartSpec[];
  error?: ErrorResponse;
}

/** AnalysisWorkbench 仅依赖 AnalysisSDK，负责把 SDK 结果转换为可展示状态。 */
export class AnalysisWorkbench {
  private state: AnalysisWorkbenchState = initialState();

  public constructor(private readonly sdk: AnalysisSDK) {}

  /** getState 返回深拷贝，避免页面代码修改工作台内部状态。 */
  public getState(): AnalysisWorkbenchState {
    return cloneState(this.state);
  }

  /** submit 执行一次分析，并稳定区分成功、空数据和 typed error。 */
  public async submit(text: string, scope: AnalysisScope): Promise<AnalysisWorkbenchState> {
    this.state = {
      status: 'loading',
      message: '正在分析…',
      charts: [],
    };

    try {
      const response = await this.sdk.analyze({ text, scope });
      this.state = {
        status: response.charts.length === 0 ? 'empty' : 'success',
        message: response.message,
        charts: cloneCharts(response.charts),
      };
    } catch (error: unknown) {
      const details = normalizeError(error);
      this.state = {
        status: 'error',
        message: details.message,
        charts: [],
        error: details,
      };
    }

    return this.getState();
  }

  /** reset 将工作台恢复为无结果的初始状态。 */
  public reset(): AnalysisWorkbenchState {
    this.state = initialState();
    return this.getState();
  }
}

function initialState(): AnalysisWorkbenchState {
  return {
    status: 'idle',
    message: '',
    charts: [],
  };
}

function normalizeError(error: unknown): ErrorResponse {
  if (error instanceof AnalysisSDKError) {
    return { ...error.details };
  }
  return {
    code: ErrorCodes.internal,
    message: 'analysis failed unexpectedly',
    retryable: false,
    requestId: 'mock-analysis-unhandled',
  };
}

function cloneState(state: AnalysisWorkbenchState): AnalysisWorkbenchState {
  return {
    ...state,
    charts: cloneCharts(state.charts),
    error: state.error === undefined ? undefined : { ...state.error },
  };
}

function cloneCharts(charts: ChartSpec[]): ChartSpec[] {
  return charts.map((chart) => ({
    ...chart,
    timeRange: { ...chart.timeRange },
  }));
}
