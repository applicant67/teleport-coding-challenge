import path from 'path';
import webpack from 'webpack';

const projectRoot = path.resolve(__dirname);

const config: webpack.Configuration = {
  node: {
    __filename: true,
    __dirname: true,
  },
  resolve: {
    alias: {
      'client': path.resolve(projectRoot, 'src'),
    },
    extensions: ['.js', '.ts', '.tsx'],
  },
  module: {
    rules: [
      {
        test: /\.ts[x]?$/,
        exclude: /node_modules/,
        use: [
          {
            loader: 'babel-loader',
            options: {
              cacheDirectory: true,
            },
          },
        ],
      },
    ],
  },
};

export default config;
