import React from 'react';
import ReactDOM from 'react-dom';
import {createStore, applyMiddleware, combineReducers} from 'redux'
import {Provider} from 'react-redux'
import thunkMiddleware from 'redux-thunk'
import reducer from './reducers'
import App from './containers/App';

console.log("1:", reducer)
let store = createStore(reducer, applyMiddleware(thunkMiddleware))
ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('root')
);
