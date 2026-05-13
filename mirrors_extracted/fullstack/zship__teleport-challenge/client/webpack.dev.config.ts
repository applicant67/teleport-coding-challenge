import path from 'path';
import webpack from 'webpack';

import base from './webpack.base.config';

const projectRoot = __dirname;

const config: webpack.Configuration = {
  ...base,
  target: 'web',
  mode: 'development',
  entry: {
    index: './src/index.ts',
  },
  devtool: 'source-map',
  output: {
    filename: '[name].js',
    path: path.resolve(projectRoot, 'build'),
  },
};

export default config;
