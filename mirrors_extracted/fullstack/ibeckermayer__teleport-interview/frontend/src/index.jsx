import React from 'react';
import ReactDOM from 'react-dom';
import {
  BrowserRouter as Router,
  Route,
  Switch,
  Redirect,
} from 'react-router-dom';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import Authenticated from './components/Authenticated';
import './index.css';
import { AppContext } from './store';

const App = () => {
  return (
    <AppContext>
      <Router>
        <Switch>
          <Route exact path="/">
            <Redirect to="/dashboard" />
          </Route>
          <Route path="/login">
            <Login />
          </Route>
          <Route path="/dashboard">
            <Authenticated>
              <Dashboard />
            </Authenticated>
          </Route>
        </Switch>
      </Router>
    </AppContext>
  );
};

ReactDOM.render(<App />, document.getElementById('root'));
