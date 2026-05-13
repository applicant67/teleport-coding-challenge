import React, { useContext } from 'react';
import { Redirect } from 'react-router-dom';
import { StoreContext } from '../../store';

// Authenticated can be used to wrap top level route components that require authentication.
// If a user attempts to access (or is in the process of accessing) an Authenticated protected
// component without a locally stored sessionID, they will be redirected to the login page.
// eslint-disable-next-line react/prop-types
const Authenticated = ({ children }) => {
  const { store } = useContext(StoreContext);

  return !store || !store.sessionID ? (
    <Redirect to="/login" />
  ) : (
    <>{children}</>
  );
};

export default Authenticated;
