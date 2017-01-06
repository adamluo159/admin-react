import React from 'react';
import ReactDOM from 'react-dom';

import Layout from './components/layout/layout'

export default class App extends React.Component {
  constructor(props){
    super(props);
  }

  render(){
    return (
      <div>
        <Layout>
        </Layout>
      </div>
    )
  }
}