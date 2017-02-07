import React from 'react';
import ReactDOM from 'react-dom';
import {
  Select,
  Message,
  Table,
  Icon,
  Button,
  Input
} from 'antd';

import {trim, checkIpFormat} from '../../utils/utils'
import {typeOption, machineColumns} from '../../constant'
const Option = Select.Option;
export default class machineMgr extends React.Component {
  constructor(props) {
    super(props);
    this.checkFunc = {
      "IP": checkIpFormat,
      "outIP": checkIpFormat
    }
  }
  componentWillMount() {
    this
      .props
      .dispatch
      .resetMachineState({editState: false})
    this
      .props
      .dispatch
      .fetchMachines();
  }

  addClick(e) {
    const {editState} = this.props.data
    if (editState) {
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return
    }

    this.editInput = {
      hostname: 'hosttmp',
      IP: "",
      outIP: "",
      type: "login"
    }
    this.new = true
    this
      .props
      .dispatch
      .addMachine({
        ...this.editInput,
        edit: true,
        key: ""
      });
  }

  SaveDo(index) {
    const {editState, data} = this.props.data
    if (!editState) {
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return
    }
    let playload = {
      index: index,
      machine:this.editInput
    }
    if (this.new) {
      this.props.dispatch.fetchAddMachine(playload)
    } else {
      this.props.dispatch.fetchSaveMachine(playload)
    }
    this.new = false
  }

  editDo(index) {
    const {editState} = this.props.data
    if (editState) {
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return
    }
    const {data} = this.props.data
    this.editInput = {
      ...data[index]
    }
    this.props.dispatch.editMachine(index)
  }

  deleteDo(index) {
    const {editState} = this.props.data
    if (editState) {
      Message.error("存在正在编辑的选项，请保存后再删除选项")
      return
    }
    this
      .props
      .dispatch
      .delMachine(index)
  }

  typeSelect(key, text, record, index) {
    let {current, pageSize} = this.props.data.page
    let {data} = this.props.data
    if (current > 1) {
      index = pageSize * (current - 1) + index
    }
    return (record.edit
      ? <div>
          <Select
            defaultValue={typeOption[0]}
            size='small'
            onChange={(value) => {
            this.editInput[key] = value
          }}>
            {typeOption.map((k) => (
              <Option value={k} key={k}>{k}</Option>
            ))
}
          </Select>
        </div>
      : <div className="editable-row-text">
        {text}
      </div>)
  }

  actionHandle(text, record, index) {
    let {current, pageSize} = this.props.data.page
    let {data} = this.props.data
    if (current > 1) {
      index = pageSize * (current - 1) + index
    }
    return (
      <div>
        {record.edit
          ? <div>
              <a onClick={(e) => {
                this.SaveDo(index)
              }}>save</a>
              <span className="ant-divider"/>
              <a onClick={(e) => {
                this.deleteDo(index)
              }}>delete</a>
            </div>
          : <div>
            <a onClick={(e) => {
              this.editDo(index)
            }}>edit</a>
            <span className="ant-divider"/>
            <a onClick={(e) => {
              this.deleteDo(index)
            }}>delete</a>
          </div>
}
      </div>
    )
  }

  actionClick(key, text, record, index) {
    let {current, pageSize} = this.props.data.page
    let {data} = this.props.data
    if (current > 1) {
      index = pageSize * (current - 1) + index
      record = data[index]
      text = data[index][key]
    }
    return (
      <div>
        {record.edit
          ? <Input
              defaultValue={text}
              size="small"
              onChange={(e) => {
              this.editInput[key] = e.target.value
            }}/>
          : <div className="editable-row-text">
            {text}
          </div>
}
      </div>
    )
  }

  render() {
    const {data, page} = this.props.data
    machineColumns.forEach((k) => {
      switch (k.key) {
        case "action":
          k.render = (text, record, index) => (this.actionHandle(text, record, index))
          break
        case "type":
          k.render = (text, record, index) => (this.typeSelect(k.key, text, record, index))
          break
        case 'describe':
          break
        default:
          k.render = (text, record, index) => (this.actionClick(k.key, text, record, index))
          break
      }
    })
    return (
      <div>
        <Button type="primary" onClick={(e) => (this.addClick(e))}>Add</Button>
        <Table
          dataSource={data}
          columns={machineColumns}
          pagination={page}
          onChange={(pagination, filters, sorter) => {
          this
            .props
            .dispatch
            .pageMachine({
              page: {
                current: pagination.current,
                pageSize: pagination.pageSize
              }
            })
        }}/>
      </div>
    )
  }
}