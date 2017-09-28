import React, { Component } from 'react'
import { Select, Message, Button, Input, Row, Col, Form, Switch, DatePicker } from 'antd'
import ZoneHead from './ZoneHead'
import { zoneConfig, zoneOptions, formItemLayout } from '../../utils/constant'

const Option = Select.Option
const FormItem = Form.Item
class ZoneFormClass extends React.Component {
  constructor(props) {
    super(props);

    const { getFieldDecorator } = this.props.form
    getFieldDecorator("edit")
    this.renderTabs(true)
  }

  dCreator(item, tag) {
    const { getFieldDecorator, getFieldsValue } = this.props.form
    let layout = item.layout ? { ...item.layout } : { ...formItemLayout }
    let options = { ...zoneOptions }
    return (
      <Col span={24} key={item.label}>
        <FormItem {...layout} label={item.label}>
          {getFieldDecorator(item.Id, options)(tag)}
        </FormItem>
      </Col>
    )
  }

  renderTabs(disabled) {
    const { channels, zoneInput, whitelst, switchEdit, datePick } = zoneConfig
    let w = this.props.form.getFieldValue("whitelst") || false
    this.renderItems = []
    for (let i = 0; i < zoneInput.length; i++) {
      this.renderItems.push(this.dCreator(zoneInput[i], <Input disabled={disabled} />))
    }
    this.renderItems.push(this.dCreator(datePick, <DatePicker disabled = {disabled} showTime format="YYYY-MM-DD HH:mm:ss" />))

    let kinds = this.props.getZoneChannels()
    let ckinds = kinds.map(k => <Option key={k}>{k}</Option>)
    this.renderItems.push(this.dCreator(channels, <Select mode={'multiple'} disabled={disabled}>{ckinds}</Select>))
    this.renderItems.push(this.dCreator(whitelst, <Switch checked ={w} disabled = {disabled} checkedChildren={'开'} unCheckedChildren={'关'} />))
  }

  render() {
    let { getFieldValue, getFieldsValue } = this.props.form
    let edit = getFieldValue("edit")
    edit ? this.renderTabs(false) : this.renderTabs(true)
    return (
      <Form onSubmit={(k) => this.handleChange(k)}>
        {this.renderItems}
      </Form>
    )
  }
}

const zoneForm = Form.create()(ZoneFormClass);
export default zoneForm
