import React, { Component } from 'react'
import { Select, Message, Button, Input, Row, Col, Form, Switch } from 'antd'
import './zone.css'
import ZoneHead from './zoneHead'
import { zoneConfig, zoneOptions, formItemLayout } from '../../constant'

const Option = Select.Option
const FormItem = Form.Item
const newZone = Form.create()(React.createClass({
  componentWillMount() {
    this.initShow = false
    this.initHead = false
    this.channelData = []
    this.ZoneHeadData = {}
    this.zoneData = {}
    this.zoneNameToZid = {}
    this.synZid = 0

    let { dispatch } = this.props
    dispatch.fetchInitZones(this.InitZones)
  },

  InitZones(json) {
    this.initHead = true
    let { setFieldsValue } = this.props.form

    if (json.Result != "OK") {
      setFieldsValue({
        edit: false
      })
      return
    }

    let zData = json.Items
    for (let i = 0; i < zData.length; i++) {
      let zone = zData[i]
      let headInfo = {
        zid: zone.zid,
        zoneName: zone.zoneName
      }

      this.zoneData[zone.zid] = zone
      this.zoneNameToZid[zone.zoneName] = zone.zid
      for (let c = 0; c < zone.channels.length; c++) {
        let channel = zone.channels[c]
        if (channel === undefined) {
          continue
        }
        if (this.ZoneHeadData[channel]) {
          this.ZoneHeadData[channel].push(headInfo)
        } else {
          this.ZoneHeadData[channel] = []
          this.ZoneHeadData[channel].push(headInfo)
          this.channelData.push(channel)
        }
      }
    }
    setFieldsValue({
      edit: false
    })
  },

  ShowZoneInfo(zid) {
    zid = Number(zid)
    this.synZid = zid
    this.initShow = true
    this.adding = false
    let zone = this.zoneData[zid]
    let { setFieldsValue } = this.props.form
    let showzone = {
      ...zone,
      edit: false
    }
    setFieldsValue(showzone)
  },

  AddZoneInfo() {
    if (this.adding == true) {
      return
    }

    this.initShow = true
    this.adding = true
    let { resetFields, setFieldsValue } = this.props.form
    resetFields()
    setFieldsValue({
      "edit": true
    })
  },

  handleChange(value) {
    value.preventDefault()
    const { form, dispatch } = this.props
    form.validateFields((err, values) => {
      if (!err) {
        values.zid = Number(values.zid)
        if (this.adding) {
          dispatch.fetchAddZone({
            obj: values,
            addZone: this.addZone
          })
        } else {
          let oldzone = this.zoneData[values.zid]
          if (oldzone == undefined) {
            let oldzid = this.zoneNameToZid[values.zoneName]
            oldzone = this.zoneData[oldzid]
          }
          if (oldzone == undefined) {
            return
          }
          dispatch.fetchSaveZone({
            obj: values,
            oldZoneName: oldzone.zoneName,
            oldZid: oldzone.zid,
            saveZone: this.saveZone
          })
        }
        this.loading = true
      }
    })
  },
  synMachine(e) {
    e.preventDefault()
    const { fetchSynMachine } = this.props.dispatch
    fetchSynMachine({ zid: this.synZid, hostname: this.zoneData[this.synZid].hostname })
  },

  addZone(json) {
    this.loading = false
    let { resetFields, setFieldsValue } = this.props.form
    let zone = json.Item
    if (json.Result != "OK") {
      let a = {
        edit: false,
      }
      setFieldsValue(a)
      return
    }
    let a = {
      ...zone,
      edit: false,
    }
    this.adding = false
    this.zoneData[zone.zid] = zone
    this.zoneNameToZid[zone.zoneName] = zone.zid
    let headInfo = {
      zid: zone.zid,
      zoneName: zone.zoneName
    }
    for (let c = 0; c < zone.channels.length; c++) {
      let channel = zone.channels[c]
      if (channel === undefined) {
        continue
      }
      if (this.ZoneHeadData[channel]) {
        this.ZoneHeadData[channel].push(headInfo)
      } else {
        this.ZoneHeadData[channel] = []
        this.ZoneHeadData[channel].push(headInfo)
        this.channelData.push(channel)
      }
    }

    setFieldsValue(a)
  },

  saveZone(rsp) {
    this.loading = false
    let { resetFields, setFieldsValue } = this.props.form

    let newZone = rsp.json.Item
    let newChannels = newZone.channels

    let oldzid = rsp.oldzid
    let oldchannels = this.zoneData[oldzid].channels
    if (rsp.json.Result != "OK") {
      let a = {
        edit: true,
      }
      setFieldsValue(a)
      return
    }
    let delFunc = (obj) => obj.zid != oldzid
    //this.nameToZid[zone.name] = zone.zid
    for (let i = 0; i < oldchannels.length; i++) {
      let delC = oldchannels[i]
      let zonelst = this.ZoneHeadData[delC]
      let newlst = zonelst.filter(delFunc)
      if (newlst.length == 0) {
        delete this.ZoneHeadData[delC]
      } else {
        this.ZoneHeadData[delC] = newlst
      }
    }
    let headInfo = {
      zid: newZone.zid,
      zoneName: newZone.zoneName
    }
    for (let c = 0; c < newZone.channels.length; c++) {
      let channel = newZone.channels[c]
      let zonelst = this.ZoneHeadData[channel]
      if (zonelst) {
        let index = zonelst.indexOf(headInfo)
        if (index === -1) {
          zonelst.push(headInfo)
        } else {
          zonelst[index] = headInfo
        }
      } else {
        zonelst = []
        zonelst.push(headInfo)
        if (this.channelData.indexOf(channel) === -1) {
          this.channelData.push(channel)
        }
      }
      this.ZoneHeadData[channel] = zonelst
    }
    let newChannelData = []
    this.channelData.forEach(k => {
      if (this.ZoneHeadData[k]) {
        newChannelData.push(k)
      }
    })
    this.channelData = newChannelData
    let a = {
      ...newZone,
      edit: false,
    }
    setFieldsValue(a)

    let oldZoneName = this.zoneData[oldzid].zoneName
    if (newZone.zoneName != oldZoneName) {
      delete this.zoneNameToZid[oldZoneName]
    }
    if (newZone.zid != oldzid) {
      delete this.zoneData[oldzid]
    }
    this.zoneData[newZone.zid] = newZone
    this.zoneNameToZid[newZone.zoneName] = newZone.zid
  },

  startZone(e) {
    e.preventDefault()
    const { fetchStartZone } = this.props.dispatch
    fetchStartZone({
      obj: { zid: this.synZid, ZoneHost: this.zoneData[this.synZid].zoneHost },
      startZoneRsp: this.startZoneRsp,
    })
  },

  startZoneRsp(jsp) {
    console.log("start zone:", jsp)
  },

  stopZone(e) {
    e.preventDefault()
    const { fetchStopZone } = this.props.dispatch
    fetchStopZone({
      obj: { zid: this.synZid, ZoneHost: this.zoneData[this.synZid].zoneHost },
      stopZoneRsp: this.startZoneRsp,
    })

  },
  stopZoneRsp(jsp) {
    console.log("stop zone:", jsp)
  },


  dCreator(item, tag) {
    const { getFieldDecorator, getFieldsValue } = this.props.form
    let layout = item.layout ? { ...item.layout } : { ...formItemLayout }
    let options = item.options ? { ...item.options } : { ...zoneOptions }
    return (
      <Col span={24} key={item.label}>
        <FormItem {...layout} label={item.label}>
          {getFieldDecorator(item.Id, options)(tag)}
        </FormItem>
      </Col>
    )
  },

  renderTabs(disabled) {
    const btnWapper = {
      span: 16,
      offset: 8
    }
    const { channels, zoneInput, whitelst } = zoneConfig
    let loading = false
    let renderItems = []
    let switchEdit = {
      Id: 'edit',
      label: '',
      layout: {
        wrapperCol: {
          span: 12
        }
      },
      options: {}
    }

    renderItems.push(this.dCreator(switchEdit, <Switch disabled={this.adding} checkedChildren={'编辑'} unCheckedChildren={'查看'} />))
    for (let i = 0; i < zoneInput.length; i++) {
      renderItems.push(this.dCreator(zoneInput[i], <Input disabled={disabled} />))
    }
    let ckinds = channels.kinds.map(k => <Option key={k}>{k}</Option>)
    renderItems.push(this.dCreator(channels, <Select multiple disabled={disabled}>{ckinds}</Select>))
    renderItems.push(this.dCreator(whitelst, <Switch disabled={disabled} checkedChildren={'开'} unCheckedChildren={'关'} />))

    return renderItems
  },

  zoneContent() {
    const { getFieldValue } = this.props.form
    let disabled = getFieldValue("edit") ? false : true
    let content = this.renderTabs(disabled)
    let buttonText = this.adding ? "新增提交" : "保存"
    return (
      <Form onSubmit={this.handleChange}>
        {content.slice(0, content.length)}
        <Row>
          <Col span={5}>
            <Button type="primary" htmlType="submit" disabled={disabled} loading={this.loading}>{buttonText}</Button>
          </Col>
          <Col span={9}>
            <Button type="primary" disabled={!disabled} onClick={(e) => this.synMachine(e)} >同步机器配置</Button>
          </Col>
          <Col span={5}>
            <Button type="primary" disabled={!disabled} onClick={(e) => this.startZone(e)} >启服</Button>
          </Col>
          <Col span={5}>
            <Button type="primary" disabled={!disabled} onClick={(e) => this.stopZone(e)} >关服</Button>
          </Col>
        </Row>
      </Form>
    )
  },

  render() {
    return (
      <div>
        <Row>
          <div id="zoneHead">
            {
              this.initHead ?
                <ZoneHead channelData={this.channelData}
                  zoneData={this.ZoneHeadData}
                  showFunc={this.ShowZoneInfo}
                  addZoneFunc={this.AddZoneInfo}
                  registerFunc={(e) => this.fresh = e}>
                </ZoneHead>
                :
                <p>Loading</p>
            }
          </div>
        </Row>
        <Row>
          <div id="zoneContent">
            <Col span={8}>
              {
                this.initShow ?
                  this.zoneContent()
                  :
                  <p> 无信息</p>
              }
            </Col>
          </div>
        </Row>
      </div>
    )
  }
}))
export default newZone
