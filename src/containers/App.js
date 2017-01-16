import React,{PropTypes} from 'react';
import ReactDOM from 'react-dom';
import { connect } from 'react-redux'
import {Icon,Row, Col } from 'antd';

import Layout from '../components/layout/layout'
import MachineMgr from '../components/machineMgr/machineMgr'
//import {selectMainlayout, 
//        fetchMachines, 
//        addMachine,
//        editMachine,
//        saveMachine
//} from '../actions'

import * as machineAction from '../actions'

const mainLays = {
  "machineMgr": (e)=>(<MachineMgr machines={e.machines} initf={e.initmachines} addmachine={e.addmachine}
                       editmachine={e.editmachine} savemachine={e.savemachine} delmachine={e.delmachine}
                       pagemachine={e.pagemachine}></MachineMgr>)
}

class App extends React.Component {
  constructor(props){
    super(props);
  }
  render(){
    const {layout,machines} = this.props
    const empty = (e)=>(console.log(""))
    const mainfunc = mainLays[layout.selectKey] || empty
    return (
     <div>
        <Row>
          <Col span={20} push={4}>      
              {mainfunc(this.props)}
          </Col>
          <Col span={4} pull={20}>
            <Layout sfunc={this.props.smainlayout}></Layout>
          </Col>
        </Row>
     </div>
    )
  }
}
App.PropTypes = {
  layout: PropTypes.shape({
    selectKey: PropTypes.string.isRequired,
  }),
  mainLayouts: PropTypes.shape({
    machines: PropTypes.arrayOf({
       editState: PropTypes.bool.isRequired,
       editInput: PropTypes.shape({
          key: PropTypes.string.isRequired,
          hostname: PropTypes.string.isRequired,
          IP: PropTypes.string.isRequired,
          outIP:PropTypes.string.isRequired,
          type :PropTypes.string.isRequired,
          edit:PropTypes.bool.isRequired,
       }),
       data: PropTypes.arrayOf(PropTypes.shape({
          key: PropTypes.string.isRequired,
          hostname: PropTypes.string.isRequired,
          IP: PropTypes.string.isRequired,
          outIP:PropTypes.string.isRequired,
          type :PropTypes.string.isRequired,
          edit:PropTypes.bool.isRequired
       })),
       page: PropTypes.shape({
         current: PropTypes.number.isRequired,
         pageSize: PropTypes.number.isRequired,
       }),
       cur: PropTypes.number.isRequired,
  })})
}

const mapStateToProps = state => {
  const {layout,machines} = state;
  return {
    layout,
    machines,
  }
}

const mapDispatchToProps = dispatch => {
  return {
    smainlayout: (e) => {dispatch(machineAction.selectMainlayout(e.key))},
    initmachines: () => {dispatch(machineAction.fetchMachines())},
    addmachine: (e)  => {dispatch(machineAction.addMachine(e))},
    editmachine:(index)=>{dispatch(machineAction.editMachine(index))},
    savemachine:(index)=>{dispatch(machineAction.saveMachine(index))},
    delmachine:(index)=>{dispatch(machineAction.delMachine(index))},
    pagemachine:(page, b, c)=>{dispatch(machineAction.pageMachine(page,b,c))}
  }
}
export default connect(mapStateToProps, mapDispatchToProps)(App);

