import type { AnalysisScope } from '../../contracts/analysis';

/** ScopeDraft 表示工作台范围输入框中的可编辑字符串。 */
export interface ScopeDraft {
  datasourceUid: string;
  from: string;
  to: string;
}

export type ScopeValidationCode =
  | 'DATASOURCE_REQUIRED'
  | 'TIME_RANGE_START_REQUIRED'
  | 'TIME_RANGE_END_REQUIRED';

/** ScopeValidationIssue 表示前端可直接展示的范围必填项问题。 */
export interface ScopeValidationIssue {
  code: ScopeValidationCode;
  field: keyof ScopeDraft;
  message: string;
}

/** createDefaultScopeDraft 返回独立且确定的 MS1 默认范围。 */
export function createDefaultScopeDraft(): ScopeDraft {
  return {
    datasourceUid: 'prometheus-mock',
    from: 'now-30m',
    to: 'now',
  };
}

/** normalizeScopeDraft 只清理输入空白，不复制 B 模块的时间语义校验。 */
export function normalizeScopeDraft(draft: ScopeDraft): ScopeDraft {
  return {
    datasourceUid: draft.datasourceUid.trim(),
    from: draft.from.trim(),
    to: draft.to.trim(),
  };
}

/** validateScopeDraft 校验前端可确定的必填项，深层时间校验仍由 B 模块负责。 */
export function validateScopeDraft(draft: ScopeDraft): ScopeValidationIssue[] {
  const normalized = normalizeScopeDraft(draft);
  const issues: ScopeValidationIssue[] = [];

  if (normalized.datasourceUid === '') {
    issues.push({
      code: 'DATASOURCE_REQUIRED',
      field: 'datasourceUid',
      message: '请选择数据源',
    });
  }
  if (normalized.from === '') {
    issues.push({
      code: 'TIME_RANGE_START_REQUIRED',
      field: 'from',
      message: '请输入开始时间',
    });
  }
  if (normalized.to === '') {
    issues.push({
      code: 'TIME_RANGE_END_REQUIRED',
      field: 'to',
      message: '请输入结束时间',
    });
  }
  return issues;
}

/** toAnalysisScope 将工作台草稿转换为 A 定义的稳定 SDK 请求结构。 */
export function toAnalysisScope(draft: ScopeDraft): AnalysisScope {
  const normalized = normalizeScopeDraft(draft);
  return {
    datasourceUid: normalized.datasourceUid,
    timeRange: {
      from: normalized.from,
      to: normalized.to,
    },
  };
}

/** describeScope 生成不依赖 UI 框架的范围摘要。 */
export function describeScope(scope: AnalysisScope): string {
  return `${scope.datasourceUid} · ${scope.timeRange.from} → ${scope.timeRange.to}`;
}
