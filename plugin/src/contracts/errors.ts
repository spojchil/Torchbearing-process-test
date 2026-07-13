/** ErrorCodes 固定前后端共享的 MS1 错误码，避免调用方使用散乱字符串。 */
export const ErrorCodes = {
  invalidArgument: 'INVALID_ARGUMENT',
  invalidScope: 'INVALID_SCOPE',
  agentUnavailable: 'AGENT_UNAVAILABLE',
  metricsUnavailable: 'METRICS_UNAVAILABLE',
  noData: 'NO_DATA',
  internal: 'INTERNAL',
} as const;

/** ErrorCode 是 ErrorCodes 中所有合法错误码的联合类型。 */
export type ErrorCode = (typeof ErrorCodes)[keyof typeof ErrorCodes];

/** ErrorResponse 表示 SDK 返回给前端的稳定错误结构。 */
export interface ErrorResponse {
  code: ErrorCode;
  message: string;
  retryable: boolean;
  requestId: string;
}

/** AnalysisSDKError 将稳定错误响应包装为可捕获的 JavaScript Error。 */
export class AnalysisSDKError extends Error {
  public readonly details: ErrorResponse;

  public constructor(details: ErrorResponse) {
    super(`${details.code}: ${details.message}`);
    this.name = 'AnalysisSDKError';
    this.details = details;
  }
}
