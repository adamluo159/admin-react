import React, {PropTypes} from 'react';
import ReactDOM from 'react-dom';
import {connect} from 'react-redux'
import {Icon, Row, Col} from 'antd';

import Layout from '../components/layout/layout'
import MachineMgr from '../components/machineMgr/machineMgr'

import actionDispatchFunc from '../actions'

const mainLays = {
  "machineMgr": (e) => {
    console.log(e.machineDispatch);
    return <MachineMgr data={e.machines} dispatch={e.machineDispatch}></MachineMgr>
  }
}

class App extends React.Component {
  constructor(props) {
    super(props);
  }
  render() {
    const {layout} = this.props
    const mainfunc = mainLays[layout.selectKey]
    return (
      <div>
        <Row>
          <Col span={20} push={4}>
            {mainfunc
              ? mainfunc(this.props)
              : null}
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
  layout: PropTypes.shape({selectKey: PropTypes.string.isRequired}),
  mainLayouts: PropTypes.shape({
    machines: PropTypes.arrayOf({
      editState: PropTypes.bool.isRequired,
      editInput: PropTypes.shape({
        key: PropTypes.string.isRequired,
        hostname: PropTypes.string.isRequired,
        IP: PropTypes.string.isRequired,
        outIP: PropTypes.string.isRequired,
        type: PropTypes.string.isRequired,
        edit: PropTypes.bool.isRequired
      }),
      data: PropTypes.arrayOf(PropTypes.shape({
        key: PropTypes.string.isRequired,
        hostname: PropTypes.string.isRequired,
        IP: PropTypes.string.isRequired,
        outIP: PropTypes.string.isRequired,
        type: PropTypes.string.isRequired,
        edit: PropTypes.bool.isRequired
      })),
      page: PropTypes.shape({current: PropTypes.number.isRequired, pageSize: PropTypes.number.isRequired}),
      cur: PropTypes.number.isRequired
    })
  })
}

const mapStateToProps = state => {
  const {layout, machines} = state;
  return {layout, machines}
}

const mapDispatchToProps = dispatch => (actionDispatchFunc(dispatch))
export default connect(mapStateToProps, mapDispatchToProps)(App);
