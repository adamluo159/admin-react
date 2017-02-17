import React, { Component } from 'react'
import { Select, Message, Button, Input, Row, Col, Form, Switch } from 'antd'
import './zone.css'
import ZoneHead from './zoneHead'
import { zoneConfig, zoneOptions, formItemLayout } from '../../constant'

const Option = Select.Option
const FormItem = Form.Item
const zone = Form.create()(React.createClass({
  componentWillMount() {
    this.initShow = false
    this.initHead = false
    this.channelData = []
    this.ZoneHeadData = {}
    this.zoneData = {}

    let {dispatch} = this.props
    dispatch.fetchInitZones(this.InitZones)
  },

  InitZones(json) {
    this.initHead = true
    let {setFieldsValue} = this.props.form

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
    console.log("cccccc", json)
    setFieldsValue({
      edit: false
    })
  },

  ShowZoneInfo(zid) {
    zid = Number(zid)
    this.initShow = true
    this.adding = false
    let zone = this.zoneData[zid]
    let {setFieldsValue} = this.props.form
    let showzone = {
      ...zone,
      edit: false
    }
    console.log("show:", zone)
    setFieldsValue(showzone)
  },

  AddZoneInfo(AddResult) {
    if (this.adding == true) {
      return
    }
    console.log(AddResult)

    this.initShow = true
    this.adding = true
    this.AddResult = AddResult
    let {resetFields, setFieldsValue} = this.props.form
    resetFields()
    setFieldsValue({
      "edit": true
    })
  },

  handleChange(value) {
    value.preventDefault()
    const {form, dispatch} = this.props
    form.validateFields((err, values) => {
      if (!err) {
        values.zid = Number(values.zid)
        dispatch.fetchAddZone({
          obj: values,
          RecvZone: this.RecvZone
        })
        this.loading = true
      }
    })
  },

  RecvZone(json) {
    this.loading = false
    let {resetFields, setFieldsValue} = this.props.form
    let zone = json.Item
    if (json.Result == "OK") {
      let a = {
        ...zone,
        edit: false,
      }
      this.adding = false
      setFieldsValue(a)
    }
    this.zoneData[zone.zid] = zone
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

    this.AddResult(json.Item)
  },

  dCreator(item, tag) {
    const {getFieldDecorator, getFieldsValue} = this.props.form
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
    const {channels, zoneInput, whitelst} = zoneConfig
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
    const {getFieldValue} = this.props.form
    let disabled = getFieldValue("edit") ? false : true
    let content = this.renderTabs(disabled)
    return (
      <Form onSubmit={this.handleChange}>
        {content.slice(0, content.length)}
        <Button type="primary" htmlType="submit" disabled={disabled} loading={this.loading}>提交</Button>
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
                  addZoneFunc={this.AddZoneInfo}>
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
export default zone
