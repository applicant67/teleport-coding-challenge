import { ApiClient } from '../../api/ApiClient';
import { Made } from './made';

export const makeApi = (overrides: Partial<ApiClient> = {}): Made<ApiClient> => ({
  authenticate: jest.fn(),
  ...overrides,
});
