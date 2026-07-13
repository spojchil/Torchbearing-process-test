/** ChartType 表示前后端共享的、与具体渲染组件无关的图表类型。 */
export type ChartType = 'timeseries' | 'stat' | 'table';

/** TimeRange 表示兼容 Grafana 的时间范围。 */
export interface TimeRange {
  from: string;
  to: string;
}

/** AnalysisScope 表示用户本次分析所选的数据源和时间范围。 */
export interface AnalysisScope {
  datasourceUid: string;
  timeRange: TimeRange;
}

/** AnalysisRequest 表示前端提交给分析 SDK 的稳定请求。 */
export interface AnalysisRequest {
  text: string;
  scope: AnalysisScope;
}

/** ChartSpec 表示后端返回的渲染器无关图表定义。 */
export interface ChartSpec {
  id: string;
  title: string;
  type: ChartType;
  datasourceUid: string;
  promql: string;
  timeRange: TimeRange;
}

/** AnalysisResponse 表示 MS1 的稳定成功响应；当前阶段只能返回确定性 mock 数据。 */
export interface AnalysisResponse {
  requestId: string;
  message: string;
  charts: ChartSpec[];
  mock: true;
}
