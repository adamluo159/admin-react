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
    this.hosts = {
      "hosttmp": true
    }
  }
  componentWillMount() {
    this
      .props
      .dispatch
      .fetchInitMachines(() => this.initProp());
  }

  initProp() {
    let {data} = this.props.data
    data.forEach((v, index) => {
      this.hosts[v.hostname] = true
    })
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
      type: "login",
      key: "tmp"
    }

    this
      .props
      .dispatch
      .addMachine({
        ...this.editInput,
        edit: true
      });
  }

  SaveDo(index, record) {
    const {editState, data} = this.props.data
    if (!editState) {
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return
    }
    if (!checkIpFormat(this.editInput.IP) || !checkIpFormat(this.editInput.outIP)) {
      Message.error("保存失败，IP不合法，请重新修改")
      return
    }
    if (this.hosts[this.editInput.hostname]) {
        Message.error("主机名重复，请重新修改")
        return
    }
    console.log("saveeeeeeee", this.hosts)

    this.editInput.key = this.editInput.hostname
    let newItem = record.key !== record.hostname
    let playload = {
      index: index,
      machine: this.editInput,
      cb:(newhost, oldhost) => {
         if(oldhost) {delete this.hosts[oldhost]}
         this.hosts[newhost]= true
        }
    }

    if (newItem) {
      this
        .props
        .dispatch
        .fetchAddMachine(playload)
    } else {
      playload.oldmachine = record
      this
        .props
        .dispatch
        .fetchSaveMachine(playload)
    }
  }

  editDo(index,record) {
    const {editState} = this.props.data
    if (editState) {
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return
    }
    const {data} = this.props.data
    this.editInput = {
      ...data[index]
    }
    this
      .props
      .dispatch
      .editMachine(index)
    delete this.hosts[record.hostname] 
  }

  deleteDo(index, record) {
    const {editState} = this.props.data
    if (editState) {
      Message.error("存在正在编辑的选项，请保存后再删除选项")
      return
    }
    this
      .props
      .dispatch
      .fetchDelMachine({
        index: index,
        fetchDel: {
          hostname: record.hostname
        },
        delCB:() => delete this.hosts[record.hostname]
      })
  }

  typeSelect(text, record, index) {
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
            this.editInput[record.type] = value
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
              <a
                onClick={(e) => {
                this.SaveDo(index, record)
              }}>save</a>
              <span className="ant-divider"/>
              <a
                onClick={(e) => {
                this.deleteDo(index, record)
              }}>delete</a>
            </div>
          : <div>
            <a onClick={(e) => {
              this.editDo(index, record)
            }}>edit</a>
            <span className="ant-divider"/>
            <a onClick={(e) => {
              this.deleteDo(index, record)
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
          k.render = (text, record, index) => (this.typeSelect(text, record, index))
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