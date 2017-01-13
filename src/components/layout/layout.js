import React from 'react';
import ReactDOM from 'react-dom';
import { Menu, Icon,Row, Col } from 'antd';
import MachineMgr from '../machineMgr/machineMgr'
import './layout.css'

const SubMenu = Menu.SubMenu;
const MenuItemGroup = Menu.ItemGroup;

let menu = {
  "游戏指令": ["machineMgr", "b", "c"],
  "主题":["1", "2", "3"]
}
const mainLays = {
  "machineMgr": (<MachineMgr></MachineMgr>)
}

export default class layout extends React.Component {
  constructor(props){
    super(props);
    this.state = {
      curKey: 'machineMgr'
    }
   }

   onclickMenu(item){
     console.log(item.key)
     this.setState({
       curKey:item.key
     })
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
    let{curKey} = this.state
    return (
      <div className = "layout-top">
        <Row>
          <Col span={20} push={4}>      
            {
              mainLays[curKey]
            }
          </Col>
          <Col span={4} pull={20}> 
          <Menu 
          onClick={(e)=>this.onclickMenu(e)}
          //style={{ width: 320 }}
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