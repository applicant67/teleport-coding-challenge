const path = require('path');

module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  parserOptions: {
    tsconfigRootDir: __dirname,
    project: [path.resolve(__dirname, 'tsconfig.json')],
  },
  plugins: ['@typescript-eslint'],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:@typescript-eslint/recommended-requiring-type-checking',
  ],
  rules: {
    // be as explicit as possible
    '@typescript-eslint/explicit-function-return-type': 'warn',

    '@typescript-eslint/switch-exhaustiveness-check': 'error',
  },
};
