export interface TimeRange {
  from: string;
  to: string;
}

export interface AnalysisScope {
  datasourceUid: string;
  timeRange: TimeRange;
}

export interface AnalysisRequest {
  text: string;
  scope: AnalysisScope;
}

export interface ChartSpec {
  id: string;
  title: string;
  type: 'timeseries' | 'stat' | 'table';
  datasourceUid: string;
  promql: string;
  timeRange: TimeRange;
}

export interface AnalysisResponse {
  requestId: string;
  chart: ChartSpec;
  mock: boolean;
}
