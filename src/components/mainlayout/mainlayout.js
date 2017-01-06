import React from 'react';
import ReactDOM from 'react-dom';
import {Select, Message, Table, Icon, Button, Input} from 'antd';
import {trim, checkIpFormat} from '../../utils/utils'
const Option = Select.Option;

const columns = [{
        title: 'HostName',
        dataIndex: 'hostname',
        key: 'hostname',
        width: '10%',
      },
      {
        title: '类型',
        dataIndex: 'type',
        key: 'type',
        width: '10%',
     },
      {
        title: '内网IP',
        dataIndex: 'IP',
        key: 'IP',
        width: '15%',
     },
      {
        title: '外网IP',
        dataIndex: 'outIP',
        key: 'outIP',
        width: '15%',
     },
     {
       title: '机器用途',
       dataIndex:'describe',
       key :'describe',
       width: '35%'
     },
      {
        title: 'Action',
        key: 'action',
        width: '15%',
}]

const select = ['login', 'master', 'zone', 'pay', 'db', 'logdb']

export default class mainlayout extends React.Component {
  constructor(props){
    super(props);
    this.state = {
      cur: 0,
      editState: false,
      editInput:[],
      data : [],
      page:{
        current: 1,
        pageSize:3,
      },
   }
   this.checkFunc={
     "IP": checkIpFormat,
     "outIP": checkIpFormat,
   }
  }

  addClick(e){
    if(this.state.editState){
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return 
    }

    let {data, columns, editInput, cur} = this.state

    let newItem = {
        key: "host" + cur,
        hostname: 'host' + cur,
        IP: "",
        outIP:"",
        type :"login",
        edit: true,
    }

    cur++
    this.setState({
      editState : true,
      data: [...data, newItem],
      editInput: [...editInput, newItem],
      cur:cur,
    });
  }

  SaveDo(index){
    if(!this.state.editState){
      Message.error("不存在正在编辑的选项")
      return 
    }

    let {data} = this.state
    let editInput = this.state.editInput[index]

    for (var k of Object.keys(this.checkFunc)) {   
      let v = editInput[k]
      if(editInput[k] !== undefined){
        if(this.checkFunc[k](v)=== false){
           Message.error("格式错误")
           return
        }
      }
    }

    editInput.edit = false
    this.setState(
      {
        data: [...data.slice(0,index), editInput, ...data.slice(index+1)],
        editState: false
      }
    )
  }

  editDo(index){
    if(this.state.editState){
      Message.error("存在正在编辑的选项，请保存后再编辑新选项")
      return 
    }
    const  {data} = this.state
    data[index].edit = true
    this.setState(
      {
        data,
        editState: true
      }
    )
  }

  deleteDo(index){
    if(this.state.editState){
      Message.error("存在正在编辑的选项，请保存后再删除选项")
      return 
    }
    let {data, editInput, page} = this.state
    let add =0
    if(index%page.pageSize==0 && page.current > 1){
      add = -1
    }
    this.setState(
      {
        data: [...data.slice(0,index), ...data.slice(index+1)],
        editInput:[...editInput.slice(0,index), ...editInput.slice(index+1)],
        page:{
          pageSize: page.pageSize,
          current:  page.current+add
        }
     }
    )
  }

  typeSelect(key, text, record, index){
    let {current, pageSize} = this.state.page
    let {data} = this.state
    if(current > 1){
      index = pageSize*(current-1) + index
    }
    return (
      record.edit?
      <div>
        <Select defaultValue={select[0]} size='small' onChange={(value)=>{this.state.editInput[index][key]= value}}>
        {
          select.map((k) => (
            <Option value={k} key={k}>{k}</Option>
          ))
        }
       </Select>
     </div>
     :
     <div className="editable-row-text">
        {text}
     </div>
    )
  }

  actionHandle(text, record, index){
    let {current, pageSize} = this.state.page
    let {data} = this.state
    if(current > 1){
      index = pageSize*(current-1) + index
    }
    return (
      <div>
      {
        record.edit?
        <div>
        <a onClick={(e)=>{this.SaveDo(index)}}>save</a>
        <span className="ant-divider" />
        <a onClick={(e)=>{this.deleteDo(index)}}>delete</a>
        </div>
        :
        <div>
        <a onClick={(e)=>{this.editDo(index)}}>edit</a>
        <span className="ant-divider" />
        <a onClick={(e)=>{this.deleteDo(index)}}>delete</a>
        </div>
      }
      </div>
    )
  }
 
  actionClick(key, text, record, index){
    let {current, pageSize} = this.state.page
    let {data} = this.state
    if(current > 1){
      index = pageSize*(current-1) + index
      record= data[index]
      text  = data[index][key]
    }
   return (
     <div>
        {
          record.edit?
          <Input defaultValue={text} size="small" onChange={(e)=>{this.state.editInput[index][key]= e.target.value}}/>
          :
          <div className="editable-row-text">
              {text}
          </div>
        }
     </div>
   )
  }

  render(){
    const {data,page} = this.state;
      
    columns.forEach((k)=>{
      switch(k.key){
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
        <Table dataSource={data} columns={columns} pagination={page} onChange={(pagination, filters, sorter)=>{
          this.setState({
            page:{
              current: pagination.current,
              pageSize:pagination.pageSize,
            }
          })
       }}/>
      </div>
    )
  }
}