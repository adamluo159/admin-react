import React from 'react';
import ReactDOM from 'react-dom';
import { Menu, Icon,Row, Col } from 'antd';
const SubMenu = Menu.SubMenu;
const MenuItemGroup = Menu.ItemGroup;
import './layout.css'

let menu = {
  "游戏指令": ["GM", "b", "c"],
  "主题":["1", "2", "3"]
}

export default class App extends React.Component {
  constructor(props){
    super(props);
   }

   onclickMenu(item){
      //menu.map((k) => {
      //       console.log(k)
      //})


     console.log(item.key)
   }

   render() {
      let parseSubIndex = (ay)=>{
             let arr = []
              ay.forEach((index) => {
                arr.push(<Menu.Item key= {index}>{index}</Menu.Item>)
              })
              return arr
      }
     let parseSubMenu = (k, ay) => (<SubMenu key={k} title={k}>{parseSubIndex(ay)}</SubMenu>)
     let parseMenu = (obj)=>{
       let arr = []
            Object.keys(obj).forEach(item => {
              arr.push(parseSubMenu(item, menu[item]))
            })
       return arr
     }

    return (
      <div className = "layout-top">
        <Row>
          <Col span={18} push={6}>      
              <div className="layout-main">
              </div>
          </Col>

          <Col span={6} pull={18}> 
          <Menu 
          onClick={this.onclickMenu}
          style={{ width: 240 }}
          //defaultOpenKeys={['sub1']}
          mode="inline"
          theme="dark">
          {parseMenu(menu)}
          </Menu>
        </Col>
        </Row>
      </div>
   )
  }
}

