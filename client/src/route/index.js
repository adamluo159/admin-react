import React, { Component } from 'react';
import { Route, Switch } from 'react-router-dom';

import Layout from '../views/Layout';
import Login from '../views/Login';
import Register from '../views/Register';

import Home from '@/views/Home';
import Machine from '@/views/MachineTable';
import Zone from '@/views/Zone';

export const childRoutes = [
  {
    'path': '/home',
    'component': Home,
    'exactly': true
  },
  {
    'path': '/machine',
    'component': Machine
  },
  {
    'path': '/zone',
    'component': Zone
  }
];

const routes = (
  <Switch>
    <Route path="/login" component={Login}/>
    <Route path="/register" component={Register}/>
    <Route path="/" component={Layout}/>
  </Switch>
);

export default routes
