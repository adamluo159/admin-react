import React from 'react';
import ReactDOM from 'react-dom';
import {connect} from 'react-redux'
import {Row, Col} from 'antd';
import Types from 'prop-types';

import Layout from '../components/layout/layout'
import MachineMgr from '../components/machineMgr/machineMgr'
import Zone from '../components/zone/zone' 

import actionDispatchFunc from '../actions'

const mainLays = {
  "machineMgr": (e) => <MachineMgr data={e.machines} dispatch={e.machineD}></MachineMgr>,
  "zone": (e) => <Zone data={e.zone} dispatch={e.zoneD}> </Zone>
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
            <Layout sfunc={this.props.layoutD.selectMainlayout}></Layout>
          </Col>
        </Row>
      </div>
    )
  }
}

App.propTypes = {
  layout: Types.shape({selectKey: Types.string}),
  mainLayouts: Types.shape({
    machines: Types.arrayOf({
      editState: Types.bool.isRequired,
      data: Types.arrayOf(Types.shape({
        key: Types.string.isRequired,
        hostname: Types.string.isRequired,
        IP: Types.string.isRequired,
        outIP: Types.string.isRequired,
        type: Types.string.isRequired,
      })),
    })
  })
}

const mapStateToProps = state => {
  const {layout, machines, zone} = state;
  return {layout, machines, zone}
}

const mapDispatchToProps = dispatch => (actionDispatchFunc(dispatch))
export default connect(mapStateToProps, mapDispatchToProps)(App);
