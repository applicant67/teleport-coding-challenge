import React from 'react';
import { render, waitFor } from '@testing-library/react';
import { screen } from '@testing-library/dom';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';
import { LoginPage } from './LoginPage';
import { mockApiWith } from '../testHelpers/mockApiWith';
import { makeApi } from '../testHelpers/makers/makeApi';
import { ApiClient } from '../api/ApiClient';
import { Made } from '../testHelpers/makers/made';
import { Routes } from '../Routes';

const mockHistory = { push: jest.fn() };
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => mockHistory,
}));

describe('LoginPage', () => {
  let mockApi: Made<ApiClient>;

  beforeEach(() => {
    mockHistory.push.mockReset();
    mockApi = makeApi({
      authenticate: jest.fn().mockReturnValue(Promise.resolve(false)),
    });
    mockApiWith(mockApi);
  });

  it('submits credentials to the api', () => {
    render(<LoginPage />);

    const username = screen.getByLabelText('Username');
    userEvent.type(username, 'some-username');

    const password = screen.getByLabelText('Password');
    userEvent.type(password, 'some-password');

    const loginButton = screen.getByText('Login');
    expect(mockApi.authenticate).not.toHaveBeenCalled();
    userEvent.click(loginButton);
    expect(mockApi.authenticate).toHaveBeenCalledWith('some-username', 'some-password');
  });

  it('shows a loading indicator and disables submit when logging in', () => {
    const loginPromise = Promise.resolve(true);
    mockApi.authenticate.mockReturnValue(loginPromise);

    render(<LoginPage />);

    userEvent.type(screen.getByLabelText('Username'), 'some-username');
    userEvent.type(screen.getByLabelText('Password'), 'some-password');

    const loginButton = screen.getByText('Login') as HTMLInputElement;
    expect(screen.queryByText('Loading...')).toBeFalsy();
    expect(loginButton.disabled).toBeFalsy();

    userEvent.click(loginButton);
    expect(screen.queryByText('Loading...')).toBeTruthy();
    expect(loginButton.disabled).toBeTruthy();
  });

  it('redirects when login succeeds', async () => {
    const loginPromise = Promise.resolve(true);
    mockApi.authenticate.mockReturnValue(loginPromise);

    render(<LoginPage />);

    userEvent.type(screen.getByLabelText('Username'), 'some-username');
    userEvent.type(screen.getByLabelText('Password'), 'some-password');

    await act(async () => {
      userEvent.click(screen.getByText('Login'));
      await loginPromise;
    });

    expect(mockHistory.push).toHaveBeenCalledWith(Routes.AUTHENTICATED);
  });

  it('shows an error when login is not successful', async () => {
    const error = new Error('Invalid credentials');
    const loginPromise = Promise.reject(error);
    mockApi.authenticate.mockReturnValue(loginPromise);

    render(<LoginPage />);

    userEvent.type(screen.getByLabelText('Username'), 'some-username');
    userEvent.type(screen.getByLabelText('Password'), 'some-password');
    const loginButton = screen.getByText('Login') as HTMLInputElement;

    expect(screen.queryByText(error.message)).toBeFalsy();
    act(() => {
      userEvent.click(loginButton);
    });
    await waitFor(() => expect(screen.queryByText(error.message)).toBeTruthy());
    expect(loginButton.disabled).toBeFalsy();
    expect(screen.queryByText('Loading...')).toBeFalsy();
    act(() => {
      userEvent.click(loginButton);
    });
    await waitFor(() => expect(screen.queryByText(error.message)).toBeFalsy());
  });

  it('requires a username and password to login', () => {
    render(<LoginPage />);
    const loginButton = screen.getByText('Login') as HTMLInputElement;
    expect(loginButton.disabled).toBeTruthy();
    userEvent.type(screen.getByLabelText('Username'), 'some-username');
    expect(loginButton.disabled).toBeTruthy();
    userEvent.type(screen.getByLabelText('Password'), 'some-password');
    expect(loginButton.disabled).toBeFalsy();
  });
});
