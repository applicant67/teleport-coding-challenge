import React, { FormEvent, useState } from 'react';
import { useHistory } from 'react-router-dom';
import { useApi } from '../api/ApiContextProvider';
import { Routes } from '../Routes';

export const LoginPage = () => {
  const history = useHistory();
  const api = useApi();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    setError('');
    setIsLoading(true);
    api.authenticate(username, password)
      .finally(() => setIsLoading(false))
      .then(() => history.push(Routes.AUTHENTICATED))
      // TODO: Error display and handling should be enhanced based on status code
      .catch((e) => setError(e.message));
  };

  const disableLogin = username === '' || password === '' || isLoading;

  return (
    <div>
      <p>
        Please provide a username and password to login.
      </p>
      {error !== '' && error}
      <form onSubmit={handleSubmit}>
        <label htmlFor="username">
          Username
          <input
            name="username"
            id="username"
            type="text"
            onChange={(event) => setUsername(event.target.value)}
          />
        </label>
        <label htmlFor="password">
          Password
          <input
            name="password"
            id="password"
            type="password"
            onChange={(event) => setPassword(event.target.value)}
          />
        </label>
        <input
          type="submit"
          disabled={disableLogin}
          value={isLoading ? 'Loading...' : 'Login'}
        />
      </form>
    </div>
  );
};
