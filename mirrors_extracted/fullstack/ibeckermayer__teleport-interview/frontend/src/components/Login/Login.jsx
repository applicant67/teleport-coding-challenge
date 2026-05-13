import React, { useContext } from 'react';
import { Redirect } from 'react-router-dom';
import api from '../../api';
import { StoreContext } from '../../store';

const Login = () => {
  const { store, setStore } = useContext(StoreContext);

  const tryLogin = async e => {
    // Form validation is handled by html5
    e.preventDefault();
    try {
      const newStore = await api.post('/login', {
        email: e.target.email.value,
        password: e.target.password.value,
      });
      setStore(newStore);
    } catch (error) {
      // TODO: check for Unauthorized and alert user that username/pwd is incorrect, remove console error
      // eslint-disable-next-line no-console
      console.error(error);
    }
  };

  return store && store.sessionID ? (
    <Redirect to="/dashboard" />
  ) : (
    <form className="login-form" onSubmit={tryLogin}>
      <h1>Sign Into Your Account</h1>

      <div>
        <label htmlFor="email">Email Address</label>
        <input
          type="email"
          id="email"
          className="field"
          autoComplete="username"
          name="email"
          required
        />
      </div>

      <div>
        <label htmlFor="password">Password</label>
        <input
          type="password"
          id="password"
          className="field"
          autoComplete="current-password"
          name="password"
          required
        />
      </div>

      <input
        type="submit"
        value="Login to my Dashboard"
        className="button block"
      />
    </form>
  );
};

export default Login;
