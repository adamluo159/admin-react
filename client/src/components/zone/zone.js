import React, { Component } from 'react'
import { Select, Message, Button, Input, Row, Col, Form, Switch } from 'antd'
import './zone.css'
import ZoneHead from './zoneHead'
import ZoneForm from './zoneForm'
import ZoneFooter from './zoneFooter'
import ZoneShowTable from './zoneShowTable'

import { zoneConfig, zoneOptions, formItemLayout } from '../../constant'

const Option = Select.Option
const FormItem = Form.Item
class ZoneClass extends React.Component {
  constructor(props) {
    super(props);
    this.zoneData = {}
    this.opZid = 0
  }

  componentWillMount() {
    let { dispatch } = this.props
    dispatch.fetchInitZones((json) => this.InitZones(json))
  }

  InitZones(json) {
    if (json.Result != "OK") {
      return
    }
    if (json.Items.length <= 0) {
      let { DisableEdit } = this.props.dispatch
      DisableEdit({ zoneEdit: false })
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

    let { DisableEdit } = this.props.dispatch
    DisableEdit({ zoneEdit: true })
    let { resetFields, setFieldsValue } = this.refs.zForm
    resetFields()
    setFieldsValue({
      "edit": true
    })

    this.refs.zFooter.setState({ addZone: true, edit: true })
  }

  saveOrAddZone(value) {
    value.preventDefault()
    const { fetchAddZone, fetchSaveZone } = this.props.dispatch
    const { getFieldsValue, setFieldsValue } = this.refs.zForm
    let zone = getFieldsValue()
    zone.zid = Number(zone.zid)
    let { addZone } = this.refs.zFooter.state
    if (addZone) {
      fetchAddZone({
        obj: zone,
        cb: (json) => this.addZoneRsp(json)
      })
    } else {
      let oldzone = this.zoneData[this.opZid]
      fetchSaveZone({
        obj: zone,
        oldZoneName: oldzone.zoneName,
        oldZid: oldzone.zid,
        cb: (json) => this.saveZoneRsp(json)
      })
    }
    setFieldsValue({ edit: false })
    this.refs.zFooter.setState({ edit: false, addZoneLoading: true })
  }

  synMachine(e) {
    e.preventDefault()
    const { fetchSynMachine } = this.props.dispatch
    const { getFieldValue } = this.refs.zForm
    let zid = Number(getFieldValue("zid"))
    fetchSynMachine({ zid: zid, hostname: this.zoneData[zid].zoneHost, cb: (json) => this.NotifyRsp(json) })
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
    const { fetchStartZone } = this.props.dispatch
    const { getFieldValue } = this.refs.zForm
    this.refs.zFooter.setState({ startZoneLoading: true })
    fetchStartZone({
      obj: { zid: this.opZid, Host: this.zoneData[this.opZid].zoneHost },
      startZoneRsp: (json) => {
        this.refs.zFooter.setState({ startZoneLoading: false })
        this.refs.zShowTable.setState({ show: json.Zstates })
        this.NotifyRsp(json)
      },
    })
  }

  stopZone(e) {
    e.preventDefault()
    const { fetchStopZone } = this.props.dispatch
    const { getFieldValue } = this.refs.zForm
    let zid = Number(getFieldValue("zid"))
    this.refs.zFooter.setState({ stopZoneLoading: true })
    fetchStopZone({
      obj: { Zid: zid, Host: this.zoneData[zid].zoneHost },
      stopZoneRsp: (json) => {
        this.refs.zFooter.setState({ stopZoneLoading: false })
        this.refs.zShowTable.setState({ show: json.Zstates })
        this.NotifyRsp(json)
      },
    })
  }
  NotifyRsp(jsp) {
    Message.warning(jsp.Result, 5);
  }

  deleteZone(e) {
    e.preventDefault()
    const { fetchDelZone } = this.props.dispatch
    this.refs.zFooter.setState({ delZoneLoading: true })
    let obj = {
      Zid: this.opZid,
      Host: this.zoneData[this.opZid].zoneHost
    }
    fetchDelZone({
      obj: obj,
      cb: (json) => this.deleteZoneRsp(json)
    })
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
    const { fetchUpdateZonelogdb } = this.props.dispatch
    const { getFieldValue } = this.refs.zForm
    let zid = Number(getFieldValue("zid"))

    fetchUpdateZonelogdb({
      Zid: zid,
      Host: this.zoneData[zid].zonelogdbHost,
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
    const { fetchStartAllZone } = this.props.dispatch
    fetchStartAllZone((json) => this.startAllZoneRsp(json))
  }

  startAllZoneRsp(json) {
    console.log("startAllZoneRsp:", json)
    this.NotifyRsp(json)
    this.refs.zFooter.setState({ startAllZoneLoading: false })
    this.refs.zShowTable.setState({ show: json.Zstates })
  }

  stopAllZone() {
    this.refs.zFooter.setState({ stopAllZoneLoading: true })
    const { fetchStopAllZone } = this.props.dispatch
    fetchStopAllZone((json) => this.stopAllZoneRsp(json))
  }

  stopAllZoneRsp(json) {
    console.log("stopAllZoneRsp:", json)
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
    let { zoneEdit } = this.props.data
    return (
      <div>
        <ZoneHead
          ref="zHead"
          addZoneFunc={() => this.AddZoneInfo()}
          showFunc={(zid) => this.ShowZone(zid)}>
        </ZoneHead>
        <Row>
          <Col span={8}>
            <div id="buttonp">
              <Switch checkedChildren={'编辑'} unCheckedChildren={'查看'} disabled={!zoneEdit} onChange={(e) => this.editZone(e)} />
              <ZoneForm ref="zForm"></ZoneForm>
            </div>
          </Col>
          <Col span={8} offset={4}>
            <ZoneShowTable ref="zShowTable" getZoneName={(zid) => this.getZoneName(zid)}></ZoneShowTable>
          </Col>
        </Row>
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
      </div>
    )
  }
}
export default ZoneClass