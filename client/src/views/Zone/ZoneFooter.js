import React, { Component } from 'react'
import { Row, Col, Button } from 'antd'



class ZoneFooter extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      edit: false,
      startZoneLoading: false,
      stopZoneLoading: false,
      startAllZoneLoading: false,
      stopAllZoneLoading: false,
      addZone: false,
      addZoneLoading: false,
      delZoneLoading: false,
      synMachineLoading: false,
      deleteZoneLoading: false,
      allZoneConfLoading: false,
    }
  }
  render() {
    const { stopZoneLoading, 
	    edit, 
	    startZoneLoading, 
	    startAllZoneLoading, 
	    stopAllZoneLoading, 
	    addZone, 
	    addZoneLoading, 
	    delZoneLoading, 
	    synMachineLoading, 
	    deleteZoneLoading,
   	    allZoneConfLoading } = this.state
		    
    const { saveOrAddZone, 
	    synMachine, 
	    startZone, 
	    stopZone, 
	    startAllZone, 
	    stopAllZone, 
	    deleteZone, 
	    updatelogZoneDB,
   	    allZoneConf } = this.props

    let buttonText = addZone ? "新增" : "保存"
    return (
      <div>
        <Row>
          <Col span={8}>
            <Row type="flex" justify="space-between">
              <Button type="primary" disabled={!edit} loading={addZoneLoading} onClick={(e) => saveOrAddZone(e)}>{buttonText}</Button>
              <Button type="primary" disabled={edit} loading={synMachineLoading} onClick={(e) => synMachine(e)} >同步机器</Button>
              <Button type="primary" disabled={edit} loading={startZoneLoading} onClick={(e) => startZone(e)} >启服</Button>
              <Button type="primary" disabled={edit} loading={stopZoneLoading} onClick={(e) => stopZone(e)} >关服</Button>
              <Button type="primary" disabled={edit} onClick={(e) => updatelogZoneDB(e)} >更新logdb</Button>
              <Button type="danger" disabled={edit} loading={deleteZoneLoading} onClick={(e) => deleteZone(e)} >删除</Button>
            </Row>
          </Col>
          <Col span={8} offset={4}>
            <Row type="flex" justify="space-around">
              <Button type="primary" loading={startAllZoneLoading} onClick={(e) => startAllZone(e)} >全服启动</Button>
              <Button type="primary" loading={stopAllZoneLoading} onClick={(e) => stopAllZone(e)} >全服关闭</Button>
              <Button type="primary" loading={allZoneConfLoading} onClick={(e) => allZoneConf(e)} >全服更新配置</Button>
            </Row>
          </Col>
        </Row>
        <div id="zoneSplit"></div>
        <Row>
          <Col span={8}>
            <Row type="flex" justify="space-between">
            </Row>
          </Col>
        </Row>
      </div >
    )
  }
}
export default ZoneFooter
