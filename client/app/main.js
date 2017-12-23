import React from 'react';
import { render } from 'react-dom';
import {
	HashRouter,
		Route,
		Link,
		Switch
} from 'react-router-dom';

import AppSider from './containers/AppSider.js'
import {Layout} from 'antd';
const {Header, Footer} = Layout

class App extends React.Component {
  constructor(props) {
		super(props);
  }
  render(){
	return (
<Layout>
	<Header>header</Header>
	<AppSider>right sidebar</AppSider>
	<Footer>footer</Footer>
</Layout>

)
  }
}

const routers = (
		<Router>
			<Route exact path="/" component={App}/>
		</Router>
		)

render(routers, document.getElementById('root'));
