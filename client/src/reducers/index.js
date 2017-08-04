import { combineReducers } from 'redux';
import auth from './auth';
import menu from './menu';
import machine from './machine';

const rootReducer = combineReducers({
  auth,
  menu,
  machine
});

export default rootReducer;
