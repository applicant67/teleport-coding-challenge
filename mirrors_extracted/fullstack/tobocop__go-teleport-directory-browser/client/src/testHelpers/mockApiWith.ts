import { ApiClient } from '../api/ApiClient';
import * as ApiContext from '../api/ApiContextProvider';

export const mockApiWith = (api: ApiClient): jest.SpyInstance => {
  const hookSpy = jest.spyOn(ApiContext, 'useApi');
  hookSpy.mockReturnValue(api);
  return hookSpy;
};
