import React from 'react';
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { withRouter } from 'react-router-dom';
import api from '../../api'
var qs = require('qs');
import './index.less';

import {
  Message,
  Table,
  Button,
  Input,
  Tag,
  Row,
  Col,
} from 'antd';

import { trim, checkIpFormat, checkHostName, checkAppliactionType } from '../../utils'
import { typeOption, machineColumns, Commonhost } from '../../utils/constant'

class MachineTable extends React.Component {
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
    this.state = {
      data: [],
      pagination: {
	 pageSize:18,
      },
      svnUpdate: false,
      commonlding: false,
    }
  }
  componentDidMount() {
    api.get('/machine').then((res) => this.ResponseAllMachine(res))
  }

  ResponseAllMachine(res) {
    if (res.error) {
      message.error(res.response.data.message);
    }
    if (typeof (res.data) == "string" && res.data == "verify") {
      this.props.history.replace('/login');
      return
    }

    if (!res.error && res.data) {
      this.initProp(res.data.Items)
    }
  }

  initProp(Items) {
    if (Items == null) {
      return
    }
    this.hosts = {}
    this.specialHosts = {}
    Items.forEach((element, index) => {
      if (element.applications != null) {
        element.applications = trim(element.applications.toString())
      }

      if (Commonhost[element.hostname]) {
        this.specialHosts[element.hostname] = element.IP
      }
      this.hosts[element.hostname] = index
    });

    Items.forEach((element, index) => {
      if (checkHostName(element.hostname)) {
        for (var sk in this.specialHosts) {
          if (this.specialHosts.hasOwnProperty(sk)) {
            var IP = this.specialHosts[sk];
            if (IP == element.IP) {
              if (element.applications == "") {
                element.applications = sk
              } else {
                element.applications = element.applications + "," + sk
              }
            }
          }
        }
      }
      element.edit = false
    })
    this.editState = false
    Items.sort(this.sortTable)
    this.setState({ data: Items })
    this.Items = Items
  }

  addClick(e) {
    if (this.editState) {
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return
    }
    this.editState = true

    this.editInput = {
      hostname: 'new',
      IP: "",
      outIP: "",
      key: "tmp",
      edit: true,
    }

    let { data, pagination } = this.state

    let index = 0
    if (pagination.current > 1) {
      index = pagination.pageSize * (pagination.current - 1)
    }
    console.log("aaaa-----", index, data)
    this.setState({
      data: [
	...data.slice(0, index),
        this.editInput,
	...data.slice(index+1),
      ],
      pagination: pagination
    })
  }

  SvnUpdateRsp(json, all) {
    if (typeof (json) == "string" && json == "verify") {
      this.props.history.replace('/login');
      return
    }

    if (json.Result != "OK") {
      Message.warning("svn更新失败", 5);
    } else {
      Message.warning("svn更新成功", 5);
    }
    this.setState({ svnUpdate: false })
  }

  svnUpdate(e) {
    this.setState({ svnUpdate: true })
    api.post('/machine/svnUpdate', { HostName: e.hostname }).then(res => this.SvnUpdateRsp(res.data, false))
  }
  svnUpdateAll(e) {
    this.setState({ svnUpdate: true })
    api.get('/machine/svnUpdateAll').then(res => this.SvnUpdateRsp(res.data, true))
  }

  SaveDo(index, record) {
    for (var key in this.editInput) {
      if (this.editInput.hasOwnProperty(key)) {
        var element = this.editInput[key];
        if (typeof (element) === "string") {
          this.editInput[key] = trim(element)
        }
      }
    }
    if (!this.editState) {
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return
    }
    if (!checkIpFormat(this.editInput.IP) || !checkIpFormat(this.editInput.outIP)) {
      Message.error("保存失败，IP不合法，请重新修改")
      return
    }
    if (this.hosts[this.editInput.hostname] != null) {
      Message.error("主机名重复，请重新修改")
      return
    }
    if (!checkHostName(this.editInput.hostname) && !Commonhost[this.editInput.hostname]) {
      Message.error("保存失败，命名不合规则，请重新修改")
      return
    }

    this.editState = true
    let newItem = record.key == "new"

    if (newItem) {
      api.post('/machine/add', qs.stringify(this.editInput)).then(res => this.ResponseAllMachine(res))
    } else {
      if (this.editInput.applications != null) {
        this.editInput.applications = this.editInput.applications.split(",")
      }
      let playload = {
        Item: JSON.stringify(this.editInput),
        Oldhost: record.hostname,
      }
      api.post('/machine/save', qs.stringify(playload)).then(res => this.ResponseAllMachine(res))
    }
  }

  editDo(index, record) {
    if (this.editState) {
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return
    }
    this.editState = true
    const { data } = this.state
    data[index].edit = true
    this.setState({
      data: data,
    })
    this.editInput = {
      ...data[index]
    }

    delete this.hosts[record.hostname]
  }

  deleteDo(index, record) {
    if (this.editState) {
      Message.error("存在正在编辑的选项，请保存后再删除选项")
      return
    }
    this.editState = true
    let playload = { hostname: record.hostname }
    api.post('/machine/del', playload).then(res => this.ResponseAllMachine(res))
  }

  actionHandle(text, record, index) {
    let { current, pageSize } = this.state.pagination
    if (current > 1) {
      index = pageSize * (current - 1) + index
    }
    let onlineColor, onlineText
    if (record.Online == true) {
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
              }}>save</a>
            <span className="ant-divider" />
            <Tag color={onlineColor}>{onlineText}</Tag>
            <span className="ant-divider" />
          </div>
          : <div>
            <a onClick={(e) => {
              this.editDo(index, record)
            }}>edit</a>
            <span className="ant-divider" />
            <a onClick={(e) => {
              this.deleteDo(index, record)
            }}>delete</a>
            <span className="ant-divider" />
            <Tag color={onlineColor}>{onlineText}</Tag>
            <span className="ant-divider" />
            <Button type="primary" loading={this.state.svnUpdate} onClick={() => (this.svnUpdate(record))}>更新svn</Button>
          </div>
        }
      </div>
    )
  }

  actionClick(key, text, record, index) {
    let { current, pageSize } = this.state.pagination
    let { data } = this.state
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
            }}
          />
          : <div className="editable-row-text">
            {text}
          </div>
        }
      </div>
    )
  }

  sortTable(a, b) {
    if (Commonhost[a.hostname] && Commonhost[b.hostname]) {
      return 0
    }
    if (!Commonhost[a.hostname] && Commonhost[b.hostname]) {
      return 1
    }
    if (Commonhost[a.hostname] && !Commonhost[b.hostname]) {
      return -1
    }
    let astring = a.hostname.split("cghost")
    let bstring = b.hostname.split("cghost")

    let anumber = astring[1]
    let bnumber = bstring[1]

    return anumber - bnumber
  }

  filtersFunc(filterArray) {
    let filterObj = {}
    let filteredData = []
    if (filterArray.length > 0) {
      filterArray.forEach(v => filterObj[Number(v)] = true)
      for (let index = 0; index < this.Items.length; index++) {
        if (this.Items[index].applications != null) {
          let arrayStr = this.Items[index].applications.split(",")
          for (let i = 0; i < arrayStr.length; i++) {
            let str = arrayStr[i]
            let ret = checkAppliactionType(str)
            if (filterObj[ret]) {
              filteredData.push(this.Items[index])
              break
            }
          }
        }
        else if (filterObj[0]) {
          filteredData.push(this.Items[index])
        }
      }
    } else {
      filteredData = this.Items.concat();
    }
    filteredData.sort(this.sortTable)
    this.setState({ data: filteredData })
  }

  commonConfig() {
    for (var key in Commonhost) {
      if (!this.hosts[key]) {
        Message.error("cannt write common config, lack " + key + " machine Info", 10)
        return
      }
    }
    this.setState({ commonlding: true })
    api.get('/machine/common').then(res => {
      if (typeof (res.data) == "string" && res.data == "verify") {
        this.props.history.replace('/login');
        return
      }
      if (res.data.Result != "OK") {
        Message.error(res.data.Result, 5)
      } else {
        Message.success(res.data.Result, 5)
      }
      this.setState({ commonlding: false })
    })
  }

  render() {
    let { data } = this.state
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
        <Row type="flex" justify="space-between">
          <Col offset={1}>
            <Button type="primary" onClick={(e) => (this.addClick(e))}>Add</Button>
          </Col>
          <Button type="primary" loading={this.state.commonlding} onClick={(e) => (this.commonConfig(e))}>生成通用服配置</Button>
          <Button type="primary" loading={this.state.svnUpdate} onClick={(e) => (this.svnUpdateAll())}>svn全机器更新</Button>
        </Row>
        <Row className ="row-machine-dis"/>
        <Table
          size="small"
          dataSource={data}
          columns={machineColumns}
          pagination={this.state.pagination}
          onChange={(pagination, filters, sorter) => {
            this.setState({ pagination: pagination })
            if (filters.applications != null) {
              this.filtersFunc(filters.applications)
            }
          }}
        />
      </div>
    )
  }
}

export default withRouter(MachineTable)
