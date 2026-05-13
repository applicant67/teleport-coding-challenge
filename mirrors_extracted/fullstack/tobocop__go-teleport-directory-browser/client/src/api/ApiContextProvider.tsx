import React, { createContext, ReactNode, useContext } from 'react';
import { ApiClient, ApiClientImpl } from './ApiClient';

const ApiContext = createContext<ApiClient | null>(null);

export const ApiContextProvider = ({ children }: { children: ReactNode }) => (
  <ApiContext.Provider value={new ApiClientImpl()}>
    {children}
  </ApiContext.Provider>
);

export const useApi = (): ApiClient => {
  const apiClient = useContext<ApiClient | null>(ApiContext);

  if (apiClient === null) {
    throw Error('Attempted to use api client when outside of the provider context');
  }

  return apiClient;
};
