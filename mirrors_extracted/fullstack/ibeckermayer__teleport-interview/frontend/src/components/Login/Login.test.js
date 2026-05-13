/**
 * @jest-environment jsdom
 */
import React from 'react';
import { mount } from 'enzyme';
import Login from './Login';
import { StoreContext } from '../../store';
import toJson from 'enzyme-to-json';

let store = {};
let setStore = newStoreValues => {
  const newStore =
    newStoreValues === null ? null : { ...store, ...newStoreValues };
  store = newStore;
};

test('it renders consistently', () => {
  const wrapper = mount(
    <StoreContext.Provider value={{ store, setStore }}>
      <Login />
    </StoreContext.Provider>
  );

  expect(toJson(wrapper)).toMatchSnapshot();
});
