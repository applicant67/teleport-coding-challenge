/* eslint-disable import/no-extraneous-dependencies */
// Configures the unofficial adapter for React 17 for Enzyme.
import { configure } from 'enzyme';
import Adapter from '@wojtekmaj/enzyme-adapter-react-17';

configure({ adapter: new Adapter() });
