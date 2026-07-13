import type { AnalysisRequest, AnalysisResponse } from '../contracts/analysis';

/** AnalysisSDK 是由 MS1 传输层或 mock owner 实现的稳定前端调用边界。 */
export interface AnalysisSDK {
  analyze(request: AnalysisRequest): Promise<AnalysisResponse>;
}
