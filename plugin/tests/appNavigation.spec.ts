import { test, expect } from './fixtures';
import { ROUTES } from '../src/constants';

test.describe('navigating app', () => {
  test('analysis workbench should render successfully', async ({ gotoPage, page }) => {
    await gotoPage(`/${ROUTES.Analysis}`);
    await expect(page.getByRole('heading', { name: 'Prometheus 指标分析' })).toBeVisible();
    await expect(page.getByRole('button', { name: '生成临时图表' })).toBeVisible();
  });
});
