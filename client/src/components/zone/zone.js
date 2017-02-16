import React, { Component } from 'react'
import { Select, Message, Button, Input, Row, Col, Form, Switch } from 'antd'
import './zone.css'
import ZoneHead from './zoneHead'
import { zoneConfig, zoneOptions, formItemLayout, zoneData } from '../../constant'

const Option = Select.Option
const FormItem = Form.Item
const zone = Form.create()(React.createClass({
  componentWillMount() {
    this.renderItems = []
    this.init = false

    let zData = Object.keys(zoneData)
    this.channelData=[]
    this.ZoneHeadData={}
    for(let i =0; i < zData.length; i++){
      let zone = zoneData[zData[i]]
      let headInfo = {
        zid: zone.zid,
        zoneName: zone.zoneName
      }
      for(let c=0; c < zone.channels.length; c++){
        let channel = zone.channels[c]
        if (channel === undefined){
          continue
        }
        if(this.ZoneHeadData[channel]){
          this.ZoneHeadData[channel].push(headInfo)
        }else{
          this.ZoneHeadData[channel]=[]
          this.ZoneHeadData[channel].push(headInfo)
          this.channelData.push(channel)
        }
      }
    }
  },
  ShowZoneInfo(zid) {
      this.init = true
      let zone = zoneData[zid]
      let {setFieldsValue} = this.props.form
      let showzone = {
          ...zone,
          edit: false
      }
      setFieldsValue(showzone)
  },

  handleChange(value) {
    value.preventDefault()
    const {form} = this.props
    form.validateFields((err, values) => {
      if (!err) {
        console.log('Received values of form: ', values)
      }
    })
    let a = form.getFieldsValue()
    console.log('aaaa', a)
  },

  dCreator(item, tag) {
    const {getFieldDecorator, getFieldsValue} = this.props.form
    let layout = item.layout ? {...item.layout} : {...formItemLayout}
    let options = item.options ? {...item.options} : {...zoneOptions}
    console.log(item, layout)
    return (
      <Col span={24} key={item.label}>
      <FormItem {...layout} label={item.label}>
        {getFieldDecorator(item.Id, options)(tag)}
      </FormItem>
      </Col>
    )
  },

  renderTabs() {
    const btnWapper = {
      span: 16,
      offset: 8
    }
    const {channels, zoneInput, whitelst} = zoneConfig
    const {getFieldValue} = this.props.form
    let loading = false
    let disabled = getFieldValue("edit") ? false : true
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
    renderItems.push(this.dCreator(switchEdit, <Switch checkedChildren={'编辑'} unCheckedChildren={'查看'} />))
    for (let i = 0; i < zoneInput.length; i++) {
      renderItems.push(this.dCreator(zoneInput[i], <Input disabled={disabled} />))
    }
    let ckinds = channels.kinds.map(k => <Option key={k}>{k}</Option>)
    renderItems.push(this.dCreator(channels, <Select multiple disabled={disabled}>{ckinds}</Select>))
    renderItems.push(this.dCreator(whitelst, <Switch disabled={disabled} checkedChildren={'开'} unCheckedChildren={'关'} />))

    return renderItems
  },

  zoneContent() {
    let content=this.renderTabs()
    return (
      <Form onSubmit={this.handleChange}>
        {content.slice(0, content.length)}
        <Button type="primary" htmlType="submit">Submit</Button>
      </Form>
    )
  },

  render() {
    return (
      <div>
       <Row>
          <div id='leftSelect'>
            <ZoneHead channelData={this.channelData} zoneData={this.ZoneHeadData} showFunc={this.ShowZoneInfo}></ZoneHead>
          </div>
       </Row>
       <Row>
          <Col span={8}>
          {
            this.init ? 
              this.zoneContent()
            :
              <p> 无信息</p>
          }
          </Col>
       </Row>
      </div>
    )
  }
}))
export default zone
