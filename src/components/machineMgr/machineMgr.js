import React from 'react';
import ReactDOM from 'react-dom';
import {Select, Message, Table, Icon, Button, Input} from 'antd';
import {trim, checkIpFormat} from '../../utils/utils'
import {typeOption, machineColumns} from '../../constant'
const Option = Select.Option;
export default class machineMgr extends React.Component {
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
  componentWillMount(){
    console.log("machines")
    this.props.initf();
  }

  addClick(e){
    const {editState} = this.props.machines
    if(editState){
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return 
    }

    this.editInput = {
     hostname : 'hosttmp',
     IP : "",
     outIP : "",
     type : "login",
    }
    this.props.addmachine({
      ...this.editInput,
      edit:true,
      key:"",
    });
  }

  SaveDo(index){
    const {editState, data} = this.props.machines
    console.log(editState)
    if(!editState){
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return 
    }
    console.log("lalala", data, "ddd", this.editInput, index)
    let editInput = this.editInput
    this.props.savemachine({
      index,
      editInput,
    })
    //const {data} = this.props.machines

    //let {data} = this.state
    //let editInput = this.state.editInput[index]

    //for (var k of Object.keys(this.checkFunc)) {   
    //  let v = editInput[k]
    //  if(editInput[k] !== undefined){
    //    if(this.checkFunc[k](v)=== false){
    //       Message.error("格式错误")
    //       return
    //    }
    //  }
    //}
    //editInput.edit = false
    //this.setState(
    //  {
    //    data: [...data.slice(0,index), editInput, ...data.slice(index+1)],
    //    editState: false
    //  }
    //)
  }

  editDo(index){
    const {editState} = this.props.machines
    if(editState){
      Message.error("存在正在编辑的选项，请保存后再添加!")
      return 
    }
    const {data} = this.props.machines
    this.editInput = {
      ...data[index]
    }
    this.props.editmachine(index)
 }

  deleteDo(index){
    const {editState} = this.props.machines
    if(editState){
      Message.error("存在正在编辑的选项，请保存后再删除选项")
      return 
    }
    this.props.delmachine(index)
 
    //let {data, editInput, page} = this.state
    //let add =0
    //if(index%page.pageSize==0 && page.current > 1){
    //  add = -1
    //}
    //this.setState(
    //  {
    //    data: [...data.slice(0,index), ...data.slice(index+1)],
    //    editInput:[...editInput.slice(0,index), ...editInput.slice(index+1)],
    //    page:{
    //      pageSize: page.pageSize,
    //      current:  page.current+add
    //    }
    // }
    //)
  }

  typeSelect(key, text, record, index){
    let {current, pageSize} = this.props.machines.page
    let {data} = this.props.machines
    if(current > 1){
      index = pageSize*(current-1) + index
    }
    return (
      record.edit?
      <div>
        <Select defaultValue={typeOption[0]} size='small' onChange={(value)=>{this.editInput[key]= value}}>
        {
          typeOption.map((k) => (
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
    let {current, pageSize} = this.props.machines.page
    let {data} = this.props.machines
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
    let {current, pageSize} = this.props.machines.page
    let {data} = this.props.machines
    if(current > 1){
      index = pageSize*(current-1) + index
      record= data[index]
      text  = data[index][key]
    }
   return (
     <div>
        {
          record.edit?
          <Input defaultValue={text} size="small" onChange={(e)=>{this.editInput[key]= e.target.value}}/>
          :
          <div className="editable-row-text">
              {text}
          </div>
        }
     </div>
   )
  }

  render(){
    //const {data,page} = this.state;
    const {data, page} = this.props.machines
    machineColumns.forEach((k)=>{
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
        <Table dataSource={data} columns={machineColumns} pagination={page} onChange={(pagination, filters, sorter)=>{
          this.props.pagemachine(pagination, filters, sorter)
       }}/>
      </div>
    )
  }
}