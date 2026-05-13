import React, { useContext, useState } from 'react';
import api from '../../api';
import { StoreContext } from '../../store';
import { useInterval } from '../../hooks';

function Dashboard() {
  const { setStore } = useContext(StoreContext);

  // TODO: save and restore state to/from localStorage to improve UX on login
  // Currently the better UX is to start at ENTERPRISE so that the Upgrade button
  // appears slightly after load for free users (on first pollMetrics result), vs starting at FREE where the Upgrade
  // button then jarringly dissapears for enterprise users slightly after load
  const [state, setState] = useState({
    plan: 'ENTERPRISE', // TODO: make global const
    maxUsers: 1000, // TODO: make global const
    totalUsers: 0,
  });

  const [showUpgradedBanner, setShowUpgradedBanner] = useState(false);

  // Polling function to get a near-continuous stream of metrics for this account
  const pollMetrics = async () => {
    try {
      const newState = await api.get('/metrics');
      setState(newState);
    } catch (error) {
      // if session timeout (in the typical case)
      if (error.response.status === 401) {
        // setStore to null to remove the session
        // <Authenticated> component will redirect the user to login page
        setStore(null);
      } else {
        // TODO send to some error logging service, for now just log. Don't necessarily
        // want to setStore(null) here, since it might just be one bad response and polling can
        // continue. Perhaps after a certain number in a row the app should setStore(null) and bail out.
        // eslint-disable-next-line no-console
        console.error(error);
      }
    }
  };

  // Set interval to poll for metrics
  // TODO: make poll interval a config var
  useInterval(pollMetrics, 300);

  const logout = async e => {
    e.preventDefault();
    try {
      await api.delete('/logout');
      setStore(null); // Delete global store on success, <Authenticated> will handle re-routing to /login
    } catch (error) {
      // TODO: Probably should alert the user something along the lines of
      // "application error, please try again and contact customer support if the error persists"
      setStore(null);
    }
  };

  // Calculates the width for the progress bar based on state and returns
  // a string with '%' appended for use in a style attribute
  const progressWidth = () => {
    const num =
      state.totalUsers <= state.maxUsers ? state.totalUsers : state.maxUsers;
    const den = state.maxUsers;
    return `${(num / den) * 100}%`;
  };

  // Core logic for flashing the upgrade banner on screen for 3 seconds after
  // successful upgrade
  const flashUpgradeBanner = () => {
    setShowUpgradedBanner(true);
    setTimeout(() => {
      setShowUpgradedBanner(false);
    }, 3000);
  };

  const upgrade = async () => {
    try {
      const newState = await api.patch('/upgrade');
      setState(newState);
      flashUpgradeBanner();
    } catch (error) {
      setStore(null);
    }
  };

  return (
    <div>
      <header className="top-nav">
        <h1>
          <i className="material-icons">supervised_user_circle</i>
          User Management Dashboard
        </h1>
        <button className="button is-border" type="button" onClick={logout}>
          Logout
        </button>
      </header>
      {state.totalUsers > state.maxUsers ? (
        <div className="alert is-error">
          You have exceeded the maximum number of users for your account, please
          upgrade your plan to increase the limit.
        </div>
      ) : null}
      {showUpgradedBanner ? (
        <div className="alert is-success">
          Your account has been upgraded successfully!
        </div>
      ) : null}
      <div className="plan">
        {state.plan === 'FREE' ? (
          <header>Startup Plan - $100/Month</header>
        ) : (
          <header>Enterprise Plan - $1000/Month</header>
        )}

        <div className="plan-content">
          <div className="progress-bar">
            <div
              style={{ width: `${progressWidth()}` }}
              className="progress-bar-usage"
            />
          </div>

          <h3>
            Users:{' '}
            {state.totalUsers <= state.maxUsers
              ? state.totalUsers
              : state.maxUsers}
            /{state.maxUsers}
          </h3>
        </div>

        <footer>
          {state.plan === 'FREE' ? (
            <button
              className="button is-success"
              type="button"
              onClick={upgrade}
            >
              Upgrade to Enterprise Plan
            </button>
          ) : null}
        </footer>
      </div>
    </div>
  );
}

export default Dashboard;
