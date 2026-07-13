import React, { FormEvent, useState } from 'react';
import { PluginPage } from '@grafana/runtime';
import { Alert, Button, Field, TextArea } from '@grafana/ui';
import { analyze } from '../api/client';
import type { AnalysisResponse } from '../api/types';
import { testIds } from '../components/testIds';

export default function AnalysisPage() {
  const [text, setText] = useState('查看 checkout 服务过去 30 分钟的请求速率');
  const [result, setResult] = useState<AnalysisResponse>();
  const [error, setError] = useState<string>();
  const [loading, setLoading] = useState(false);

  const onSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setLoading(true);
    setError(undefined);

    try {
      setResult(
        await analyze({
          text,
          scope: {
            datasourceUid: 'prometheus-default',
            timeRange: { from: 'now-30m', to: 'now' },
          },
        })
      );
    } catch (err) {
      setError(err instanceof Error ? err.message : '分析请求失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <PluginPage>
      <main data-testid={testIds.analysis.container}>
        <h1>Prometheus 指标分析</h1>
        <p>MS1 空骨架：提交自然语言，后端返回一张临时图表规格的 Mock。</p>

        <form onSubmit={onSubmit}>
          <Field label="分析需求">
            <TextArea
              data-testid={testIds.analysis.input}
              value={text}
              rows={4}
              onChange={(event) => setText(event.currentTarget.value)}
            />
          </Field>
          <Button data-testid={testIds.analysis.submit} type="submit" disabled={!text.trim() || loading}>
            {loading ? '生成中…' : '生成临时图表'}
          </Button>
        </form>

        {error && <Alert title="请求失败">{error}</Alert>}
        {result && (
          <section data-testid={testIds.analysis.result}>
            <h2>{result.chart.title}</h2>
            <p>图表类型：{result.chart.type}</p>
            <p>
              查询状态：{result.query.status}，序列数：{result.query.seriesCount}，耗时：{result.query.durationMs}ms
            </p>
            <pre>{result.chart.promql}</pre>
            <p>{result.evidence.explanation}</p>
            {result.mock && <small>当前返回架构 Stub，真实 Agent、MCP transport 与查询由后续 Issue 实现。</small>}
          </section>
        )}
      </main>
    </PluginPage>
  );
}
