import {
  createDefaultScopeDraft,
  describeScope,
  normalizeScopeDraft,
  toAnalysisScope,
  validateScopeDraft,
} from './scopeModel';

// @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
describe('scopeModel', () => {
  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('提供确定性的默认范围', () => {
    const draft = createDefaultScopeDraft();
    const scope = toAnalysisScope(draft);

    assertEqual(describeScope(scope), 'prometheus-mock · now-30m → now');
  });

  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('清理输入空白但不重写 Grafana 时间表达式', () => {
    const normalized = normalizeScopeDraft({
      datasourceUid: ' prometheus-mock ',
      from: ' now/d ',
      to: ' now/d+1d ',
    });

    assertEqual(normalized.datasourceUid, 'prometheus-mock');
    assertEqual(normalized.from, 'now/d');
    assertEqual(normalized.to, 'now/d+1d');
  });

  // @ts-ignore -- Jest 全局函数由已声明但尚未安装的测试依赖提供。
  it('为三个空字段返回稳定的边界问题', () => {
    const issues = validateScopeDraft({ datasourceUid: ' ', from: '', to: ' ' });

    assertEqual(issues.length, 3);
    assertEqual(issues[0]?.code, 'DATASOURCE_REQUIRED');
    assertEqual(issues[1]?.code, 'TIME_RANGE_START_REQUIRED');
    assertEqual(issues[2]?.code, 'TIME_RANGE_END_REQUIRED');
  });
});

function assertEqual(actual: unknown, expected: unknown): void {
  if (actual !== expected) {
    throw new Error(`actual ${String(actual)} does not equal expected ${String(expected)}`);
  }
}
