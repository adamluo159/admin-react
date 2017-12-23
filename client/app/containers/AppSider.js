import React from 'react';
import {Menu,Layout} from 'antd';
const {Sider} = Layout

const Mune = [
{
	key : "1",
	    name:"aaaa",
},
{
	key : "2",
	name:"bbbb",
}
]


export default class AppSider extends React.Component {
	constructor(props) {
		super(props);
	}

	menuClick(e){
		console.log("aaaaaaaa--",e)
	}
	render() {
		return (
				<Sider>
					<Menu theme="dark" mode="inline" onClick ={this.menuClcik}>
					</Menu>
				</Sider>
				)
	}
}
