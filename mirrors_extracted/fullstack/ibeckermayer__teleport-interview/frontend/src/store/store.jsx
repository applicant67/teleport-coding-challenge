import React, { createContext, useState } from 'react';
import { getLocalStore, setLocalStore } from '../localStorage';

// Global React.Context value for components to pass into useContext in order to access the store
const StoreContext = createContext();

// AppContext should wrap the entire application, so that any component in the application
// can access the global stores's state as a StoreContext.Consumer
// eslint-disable-next-line react/prop-types
const AppContext = ({ children }) => {
  // Initialize the app with store set to the store saved in localStorage (or null if no saved store exists)
  // setStoreI is "I" for internal use only (like a private method)
  const [store, setStoreI] = useState(getLocalStore());

  // Consumers can update the store by calling setStore, or clear the store by passing in null
  const setStore = newStoreValues => {
    const newStore =
      newStoreValues === null ? null : { ...store, ...newStoreValues };
    setStoreI(newStore); // update store
    setLocalStore(newStore); // keep localStorage in sync
  };

  return (
    <StoreContext.Provider value={{ store, setStore }}>
      {children}
    </StoreContext.Provider>
  );
};

export { StoreContext, AppContext };
