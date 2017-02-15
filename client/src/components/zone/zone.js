import React, { Component } from 'react'
import { Select, Message, Button, Input, Row, Col, Form, Switch } from 'antd'
import './zone.css'
import ZoneHead from './zoneHead'
import { zoneConfig, zoneOptions, formItemLayout, ZoneData } from '../../constant'

const Option = Select.Option
const FormItem = Form.Item
const zone = Form.create()(React.createClass({
  componentWillMount() {
    this.renderItems = []
    this.init = false
    this.data = {
      ...ZoneData
    }

    this.initZone = {
      channels: ['IOS','yyb'],
      edit: false,
      whitelst: true,
      zid: "1",
      zoneDBHost: "1", 
      zoneHost: "1", 
      zoneName: "1",
      zonelogdbHost: "1"
    }
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
    
    let curzone = getFieldsValue()
    console.log(item.Id, curzone[item.Id], this.initZone[item.Id])

    if (curzone[item.Id] == undefined){
        options.initialValue=this.initZone[item.Id]
    }
    console.log("aaa", options.initialValue)
    
    return (
      <Col span={24} key={item.label}>
      <FormItem layout label={item.label}>
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
    const {channels, zoneInput, wihtelst} = zoneConfig
    const {getFieldValue} = this.props.form
    let loading = false
    let disabled = getFieldValue("edit")

    this.renderItems = []
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
    this.renderItems.push(this.dCreator(switchEdit, <Switch checkedChildren={'查看'} unCheckedChildren={'编辑'} />))

    for (let i = 0; i < zoneInput.length; i++) {
      this.renderItems.push(this.dCreator(zoneInput[i], <Input disabled={disabled} />))
    }
    let ckinds = channels.kinds.map(k => <Option key={k}>
                                           {k}
                                         </Option>)
    this.renderItems.push(this.dCreator(channels, <Select multiple disabled={disabled}>
                                                    {ckinds}
                                                  </Select>))
    this.renderItems.push(this.dCreator(wihtelst, <Switch disabled={disabled} checkedChildren={'开'} unCheckedChildren={'关'} />))
  },

  zoneContent(disabled) {
    this.renderTabs(disabled)
    return (
      <Form onSubmit={this.handleChange}>
        {this.renderItems.slice(0, this.renderItems.length)}
        <Button type="primary" htmlType="submit">Submit</Button>
      </Form>
    )
  },

  switchToEdit(e) {
    this.disabled = e
  },

  render() {
    return (
      <div>
        <Row>
          <Col span={24}>
          <div id='leftSelect'>
            <ZoneHead></ZoneHead>
          </div>
          </Col>
          <Col span={12}>
          {this.zoneContent()}
          </Col>
        </Row>
      </div>
    )
  }
}))
export default zone
