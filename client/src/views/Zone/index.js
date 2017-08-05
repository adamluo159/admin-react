import React, { Component } from 'react'
import { withRouter } from 'react-router-dom';
import api from '../../api'
var qs = require('qs');
import { Select, Message, Button, Input, Row, Col, Form, Switch, Layout } from 'antd'
import ZoneHead from './ZoneHead'
import ZoneForm from './ZoneForm'
import ZoneFooter from './ZoneFooter'
import './index.less';
import ZoneShowTable from './ZoneShowTable'

import { zoneConfig, zoneOptions, formItemLayout } from '../../utils/constant'
const { Header, Footer, Sider, Content } = Layout;

const Option = Select.Option
const FormItem = Form.Item
class ZoneClass extends React.Component {
  constructor(props) {
    super(props);
    this.zoneData = {}
    this.opZid = 0
    this.state = {
      zoneEdit: false,
    }
  }

  componentDidMount() {
    api.get('/zone').then((res) => this.InitZones(res.data))
  }

  InitZones(json) {
    if (json.Result != "OK") {
      return
    }
    if (json.Items.length <= 0) {
      this.setState({ zoneEdit: false })
    } else {
      json.Items.forEach(v => {
        this.zoneData[v.zid] = v
      })
      this.opZid = this.refs.zHead.Init(this.zoneData)
      let { setFieldsValue } = this.refs.zForm
      setFieldsValue(this.zoneData[this.opZid])
    }
    this.refs.zShowTable.setState({ show: json.Zstates })
  }

  ShowZone(zid) {
    this.opZid = Number(zid)
    let zone = this.zoneData[zid]
    let { setFieldsValue } = this.refs.zForm
    let showzone = {
      ...zone,
      edit: false
    }
    setFieldsValue(showzone)
  }

  AddZoneInfo() {
    if (this.addingZone == true) {
      return
    }

    let { resetFields, setFieldsValue } = this.refs.zForm
    resetFields()
    setFieldsValue({ "edit": true })
    this.refs.zFooter.setState({ addZone: true, edit: true })
    this.setState({ zoneEdit: true })
  }

  saveOrAddZone(value) {
    value.preventDefault()
    const { getFieldsValue, setFieldsValue } = this.refs.zForm
    let zone = getFieldsValue()
    zone.zid = Number(zone.zid)
    let { addZone } = this.refs.zFooter.state
    if (addZone) {
      api.post('/zone/add', zone).then((res) => this.addZoneRsp(res.data))
    } else {
      let oldzone = this.zoneData[this.opZid]
      let data = {
        Item: zone,
        OldZoneName: oldzone.zoneName,
        OldZid: oldzone.zid,
        cb: (json) => this.saveZoneRsp(json)
      }
      api.post('/zone/save', data).then((res) => this.addZoneRsp(res.data))
    }
    setFieldsValue({ edit: false })
    this.refs.zFooter.setState({ edit: false, addZoneLoading: true })
  }

  synMachine(e) {
    e.preventDefault()
    const { getFieldValue } = this.refs.zForm
    let zid = Number(getFieldValue("zid"))
    let hostname = this.zoneData[zid].zoneHost
    api.get("/zone/synMachine?zid=" + zid + "&hostname=" + hostname).then((res) => {
      this.refs.zFooter.setState({ edit: false, synMachineLoading: false })
      this.NotifyRsp(res.data)
    })
    this.refs.zFooter.setState({ edit: false, synMachineLoading: true })
  }

  addZoneRsp(json) {
    let { setFieldsValue } = this.refs.zForm
    let zone = json.Item
    if (json.Result != "OK") {
      setFieldsValue(this.zoneData[this.opZid])
      return
    }
    this.zoneData[zone.zid] = zone
    this.opZid = zone.zid
    let addContent = {
      ...zone,
      edit: false,
    }
    setFieldsValue(addContent)
    this.refs.zFooter.setState({ addZoneLoading: false, addZone: false, edit: false })
    this.refs.zHead.Init(this.zoneData, this.opZid)
  }

  saveZoneRsp(rsp) {
    let { setFieldsValue } = this.refs.zForm
    let newZone = rsp.json.Item
    let oldzid = rsp.oldzid
    if (rsp.json.Result != "OK") {
      setFieldsValue(this.zoneData(oldzid))
      return
    }
    if (oldzid == newZone.zid) {
      this.zoneData[oldzid] = newZone
    } else {
      delete this.zoneData[oldzid]
      this.zoneData[newZone] = newZone
    }
    this.refs.zHead.Init(this.zoneData, newZone.zid)
    this.refs.zFooter.setState({ edit: false, addZoneLoading: false })
  }

  startZone(e) {
    e.preventDefault()
    const { getFieldValue } = this.refs.zForm
    this.refs.zFooter.setState({ startZoneLoading: true })
    let obj = { zid: this.opZid, Host: this.zoneData[this.opZid].zoneHost }
    api.post('/zone/startZone', obj).then((res) => {
      this.refs.zFooter.setState({ startZoneLoading: false })
      this.refs.zShowTable.setState({ show: res.data.Zstates })
      this.NotifyRsp(res.data)
    })
  }

  stopZone(e) {
    e.preventDefault()
    const { getFieldValue } = this.refs.zForm
    let zid = Number(getFieldValue("zid"))
    let obj = { Zid: zid, Host: this.zoneData[zid].zoneHost }
    this.refs.zFooter.setState({ stopZoneLoading: true })
    api.post('/zone/stopZone', obj).then((res) => {
      this.refs.zFooter.setState({ stopZoneLoading: false })
      this.refs.zShowTable.setState({ show: res.data.Zstates })
      this.NotifyRsp(res.data)
    })
  }
  NotifyRsp(jsp) {
    Message.warning(jsp.Result, 5);
  }

  deleteZone(e) {
    e.preventDefault()
    this.refs.zFooter.setState({ delZoneLoading: true })
    let obj = {
      Zid: this.opZid,
      Host: this.zoneData[this.opZid].zoneHost
    }
    api.post('/zone/del', obj).then((res) => {
      this.refs.zFooter.setState({ deleteZoneLoading: false })
      this.deleteZoneRsp(res.data)
    })
    this.refs.zFooter.setState({ deleteZoneLoading: true })
  }

  deleteZoneRsp(json) {
    this.refs.zFooter.setState({ delZoneLoading: false })
    if (json.Result != "OK") {
      return
    }
    let zid = json.Item.zid
    let { setFieldsValue } = this.refs.zForm

    delete this.zoneData[zid]
    this.opZid = this.refs.zHead.Init(this.zoneData)

    setFieldsValue(this.zoneData[this.opZid])
  }

  updatelogZoneDB(e) {
    e.preventDefault()
    const { getFieldValue } = this.refs.zForm
    let zid = Number(getFieldValue("zid"))
    let data = {
      Zid: zid,
      Host: this.zoneData[zid].zonelogdbHost,
    }

    api.post('/zone/updateZonelogdb', obj).then((res) => {
      this.NotifyRsp(res.data)
    })
  }

  editZone(e) {
    const { setFieldsValue } = this.refs.zForm
    let edit = { edit: e }
    setFieldsValue(edit)
    this.refs.zFooter.setState(edit)
  }

  startAllZone() {
    this.refs.zFooter.setState({ startAllZoneLoading: true })
    api.post('/zone/startAllZone').then((res) => {
      this.startAllZoneRsp(res.data)
    })
  }

  startAllZoneRsp(json) {
    this.NotifyRsp(json)
    this.refs.zFooter.setState({ startAllZoneLoading: false })
    this.refs.zShowTable.setState({ show: json.Zstates })
  }

  stopAllZone() {
    this.refs.zFooter.setState({ stopAllZoneLoading: true })
    api.post('/zone/stopAllZone').then((res) => {
      this.stopAllZoneRsp(res.data)
    })
  }

  stopAllZoneRsp(json) {
    this.NotifyRsp(json)
    this.refs.zFooter.setState({ stopAllZoneLoading: false })
    this.refs.zShowTable.setState({ show: json.Zstates })
  }

  getZoneName(zid) {
    let data = this.zoneData[zid]
    if (data == null) {
      return 0
    }
    return data.zoneName
  }

  render() {
    return (
      <Layout>
        <Header className="layout-head">
          <ZoneHead
            ref="zHead"
            addZoneFunc={() => this.AddZoneInfo()}
            showFunc={(zid) => this.ShowZone(zid)}>
          </ZoneHead>
        </Header>
        <Layout>
          <Content>
            <Row className="row-dis" />
            <Row>
              <Col offset={1}>
                <Switch checkedChildren={'编辑'} unCheckedChildren={'查看'} disabled={false} onChange={(e) => this.editZone(e)} />
              </Col>
            </Row>
            <Row type="flex" justify="center" align="top">
              <ZoneForm ref="zForm"></ZoneForm>
            </Row>
          </Content>
          <Sider className="layout-head">
            <ZoneShowTable ref="zShowTable" getZoneName={(zid) => this.getZoneName(zid)}></ZoneShowTable>
          </Sider>
        </Layout>
          <Footer>
            <ZoneFooter ref="zFooter"
              synMachine={(e) => this.synMachine(e)}
              startZone={(e) => this.startZone(e)}
              stopZone={(e) => this.stopZone(e)}
              deleteZone={(e) => this.deleteZone(e)}
              updatelogZone={(e) => this.updatelogZoneDB(e)}
              saveOrAddZone={(e) => this.saveOrAddZone(e)}
              stopAllZone={(e) => this.stopAllZone(e)}
              startAllZone={(e) => this.startAllZone(e)}>
            </ZoneFooter>
          </Footer>
      </Layout>
    )
  }
}
export default withRouter(ZoneClass)