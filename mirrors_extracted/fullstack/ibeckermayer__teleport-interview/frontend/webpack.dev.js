// eslint-disable-next-line import/no-extraneous-dependencies
const { merge } = require('webpack-merge');
const fs = require('fs');
const path = require('path');
const common = require('./webpack.common.js');

module.exports = merge(common, {
  mode: 'development',
  devtool: 'inline-source-map',
  devServer: {
    publicPath: '/dist',
    historyApiFallback: true, // Enables browser routing (i.e. react-router)
    host: '0.0.0.0', // Makes devServer accessible from outside of Docker container
    proxy: {
      // See https://webpack.js.org/configuration/dev-server/#devserverproxy
      // Re-routes fetch calls to 'api/...' to 'https://backend-dev:8000/api...',
      // which is how we can properly make calls to our backend server when both this
      // dev server and the backend server are running in docker containers.
      '/api': {
        target: 'https://backend-dev:8000',
        secure: false, // Prevents nodejs server from rejecting self-signed certs
      },
    },
    https: {
      key: fs.readFileSync(path.resolve(__dirname, '../certs/localhost.key')),
      cert: fs.readFileSync(path.resolve(__dirname, '../certs/localhost.crt')),
    },
  },
});
