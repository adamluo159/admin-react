import React from 'react';
import ReactDOM from 'react-dom';
import {
  Select,
  Message,
  Table,
  Icon,
  Button,
  Input,
  Tag
} from 'antd';

import { trim, checkIpFormat } from '../../utils/utils'
import { typeOption, machineColumns } from '../../constant'
const Option = Select.Option;
export default class machineMgr extends React.Component {
  constructor(props) {
    super(props);
    //记录当前已经存在的机器名
    this.hosts = {
      "hosttmp": false
    }
    this.columnsRender = {
      "hostname": (text, record, index) => (this.actionClick("hostname", text, record, index)),
      "IP": (text, record, index) => (this.actionClick("IP", text, record, index)),
      "outIP": (text, record, index) => (this.actionClick("outIP", text, record, index)),
      "action": (text, record, index) => (this.actionHandle(text, record, index))
    }
  }
  componentWillMount() {
    this
      .props
      .dispatch
      .fetchInitMachines((json) => this.initProp(json));
  }

  initProp(rsp) {
    if (rsp.Items == null) {
      return
    }
    rsp.Items.forEach(element => {
      if (element.applications != null) {
        element.applications = element.applications.toString()
      }
      this.hosts[element.hostname] = element.Online
      delete element.Online
    });
    this.props.dispatch.InitMachines({
      data: rsp.Items,
      editState: false
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
    if (this.hosts[this.editInput.hostname] != null) {
      console.log(this.hosts, this.editInput)
      Message.error("主机名重复，请重新修改")
      return
    }
    this.editInput.key = this.editInput.hostname
    let newItem = record.key !== record.hostname
    let playload = {
      index: index,
      machine: this.editInput,
      cb: (newhost, oldhost) => {
        if (oldhost) {
          delete this.hosts[oldhost]
        }
      }
    }
    if (playload.machine.applications != null) {
      playload.machine.applications = playload.machine.applications.split(","); 
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

  editDo(index, record) {
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
        delCB: () => delete this.hosts[record.hostname]
      })
  }

  actionHandle(text, record, index) {
    let {current, pageSize} = this.props.data.page
    if (current > 1) {
      index = pageSize * (current - 1) + index
    }
    let onlineColor, onlineText
    if (this.hosts[record.hostname] == true) {
      onlineColor = "green"
      onlineText = "已连接"
    } else {
      onlineColor = "pink"
      onlineText = "未连接"
    }
    return (
      <div>
        {record.edit
          ? <div>
            <a
              onClick={(e) => {
                this.SaveDo(index, record)
              } }>save</a>
            <span className="ant-divider" />
            <Tag color={onlineColor}>{onlineText}</Tag>
          </div>
          : <div>
            <a onClick={(e) => {
              this.editDo(index, record)
            } }>edit</a>
            <span className="ant-divider" />
            <a onClick={(e) => {
              this.deleteDo(index, record)
            } }>delete</a>
            <span className="ant-divider" />
            <Tag color={onlineColor}>{onlineText}</Tag>
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
            } } />
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
      k.render = this.columnsRender[k.key]
        ? this.columnsRender[k.key]
        : k.render
    })
    data.forEach((k) => {
      k.key = k.hostname
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
          } } />
      </div>
    )
  }
}