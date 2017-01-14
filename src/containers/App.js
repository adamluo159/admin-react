import React,{PropTypes} from 'react';
import ReactDOM from 'react-dom';
import { connect } from 'react-redux'
import {Icon,Row, Col } from 'antd';

import Layout from '../components/layout/layout'
import MachineMgr from '../components/machineMgr/machineMgr'
import {selectMainlayout} from '../actions'

const mainLays = {
  //"machineMgr": (<MachineMgr>machines={mainLayouts.machines}</MachineMgr>)
  "machineMgr": (<MachineMgr></MachineMgr>)
}

class App extends React.Component {
  constructor(props){
    super(props);
  }
  render(){
    const {layout, mainlayouts} = this.props
    return (
     <div>
        <Row>
          <Col span={20} push={4}>      
              {mainLays[layout.selectKey]}
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
       key: PropTypes.string.isRequired,
       hostname: PropTypes.string.isRequired,
       IP: PropTypes.string.isRequired,
       outIP:PropTypes.string.isRequired,
       type :PropTypes.string.isRequired,
       edit:PropTypes.bool.isRequired,
    }),
  })

}

const mapStateToProps = state => {
  const {layout, mainlayouts} = state;
  return {
    layout,
    mainlayouts,
  }
}

const mapDispatchToProps = dispatch => {
  return {
    smainlayout: (e) => {dispatch(selectMainlayout(e.key))},
  }
}
export default connect(mapStateToProps, mapDispatchToProps)(App);

