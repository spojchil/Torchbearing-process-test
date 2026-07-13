export interface TimeRange {
  from: string;
  to: string;
}

export interface AnalysisScope {
  datasourceUid: string;
  timeRange: TimeRange;
  service?: string;
  environment?: string;
  namespace?: string;
  cluster?: string;
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
  unit?: string;
  legend?: string;
}

export interface QuerySummary {
  language: 'promql';
  expression: string;
  status: 'success' | 'no_data';
  durationMs: number;
  seriesCount: number;
}

export interface Evidence {
  metrics: string[];
  explanation: string;
}

export interface AnalysisResponse {
  requestId: string;
  chart: ChartSpec;
  query: QuerySummary;
  evidence: Evidence;
  mock: boolean;
}

export interface AnalysisErrorResponse {
  requestId: string;
  error: {
    code: string;
    message: string;
  };
}
