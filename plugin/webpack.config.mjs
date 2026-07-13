import path from 'node:path';
import { fileURLToPath } from 'node:url';
import CopyPlugin from 'copy-webpack-plugin';

const root = path.dirname(fileURLToPath(import.meta.url));

export default {
  context: root,
  entry: './src/module.tsx',
  externals: {
    '@grafana/data': '@grafana/data',
    '@grafana/runtime': '@grafana/runtime',
    '@grafana/ui': '@grafana/ui',
    react: 'react',
    'react-dom': 'react-dom',
  },
  mode: 'production',
  module: {
    rules: [
      {
        exclude: /node_modules/,
        test: /\\.[jt]sx?$/,
        use: 'swc-loader',
      },
    ],
  },
  output: {
    clean: true,
    filename: 'module.js',
    path: path.resolve(root, 'dist'),
  },
  plugins: [new CopyPlugin({ patterns: [{ from: 'plugin.json', to: 'plugin.json' }] })],
  resolve: {
    extensions: ['.ts', '.tsx', '.js'],
  },
};
