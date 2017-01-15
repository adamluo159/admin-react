import React from 'react';
import ReactDOM from 'react-dom';
import {Menu} from 'antd';
import {adminMenu} from '../../constant'
import './layout.css'
const SubMenu = Menu.SubMenu;
const MenuItemGroup = Menu.ItemGroup;

export default class layout extends React.Component {
  constructor(props){
    super(props);
  }
  render() {
      let parseSubIndex = (ay)=>{
             let arr = []
              ay.forEach((index) => {
                arr.push(<Menu.Item key= {index.key}>{index.text}</Menu.Item>)
              })
              return arr
      }
     let parseSubMenu = (k, ay) => (<SubMenu key={k} title={k}>{parseSubIndex(ay)}</SubMenu>)
     let parseMenu = (obj)=>{
       let arr = []
            Object.keys(obj).forEach(item => {
              arr.push(parseSubMenu(item, adminMenu[item]))
            })
       return arr
     }
    return (
     <Menu 
          onClick={(e)=>{this.props.sfunc(e)}}
          //style={{ width: 320 }}
          //defaultOpenKeys={['sub1']}
          mode="inline"
          theme="dark">
          {parseMenu(adminMenu)}
     </Menu>
   )
  }
}