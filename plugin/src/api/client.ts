import { getBackendSrv } from '@grafana/runtime';
import { lastValueFrom } from 'rxjs';
import { PLUGIN_ID } from '../constants';
import type { AnalysisRequest, AnalysisResponse } from './types';

export async function analyze(request: AnalysisRequest): Promise<AnalysisResponse> {
  const response = await lastValueFrom(
    getBackendSrv().fetch<AnalysisResponse>({
      url: `/api/plugins/${PLUGIN_ID}/resources/analysis`,
      method: 'POST',
      data: request,
    })
  );

  return response.data;
}
