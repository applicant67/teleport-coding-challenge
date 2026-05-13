import React from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { Routes } from './Routes';
import { LoginPage } from './login/LoginPage';
import { ApiContextProvider } from './api/ApiContextProvider';

function App() {
  const authSuccessful = () => <div>Auth Successful</div>;
  return (
    <BrowserRouter>
      <Switch>
        <ApiContextProvider>
          <Route exact strict path={Routes.LOGIN} component={LoginPage} />

          {/* TODO: Protect this and other authenticated routes */}
          <Route exact strict path={Routes.AUTHENTICATED} render={authSuccessful} />
        </ApiContextProvider>
      </Switch>
    </BrowserRouter>
  );
}

export default App;
