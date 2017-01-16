import { combineReducers } from 'redux'
//import {REQ_MACHINES, 
//        RECV_MACHINES,
//        ADD_MACHINE,
//        EDIT_MACHINE,
//        SAVE_MACHINE,
//        DEL_MACHINE,
//        SELECT_MAINLAYOUT_KEY,
//} from '../actions'

import * as machine from '../actions'

const initState = {
    loading : true,
}
const layout = (state ={selectKey: 'machineMgr'}, action)=>{
    switch (action.type) {
        case machine.SELECT_MAINLAYOUT_KEY:
            return {
                selectKey: action.selectKey,
            }
        default:
            return state
    }
}

const machinesInitState ={
    editState: false,
    editInput: {},
    data : [],
    page: {
        current: 1,
        pageSize:3,
    },
    cur: 0,
}

const addmachineData = (oldState, newItem) => {
    const {cur, page, data} = oldState
    newItem.key += cur
    newItem.hostname +=cur
    return {
        editState:true,
        page: page,
        cur: cur+1,
        data: [...data, newItem],
        editInput: newItem,
    }
}

const editmachineData = (oldState, index) => {
    const {data} = oldState
    let editItem = Object.assign({}, data[index])
    editItem.edit = true
    return {
        ...oldState,
        editState: true,
        editInput:{
            ...oldState[index],
        },
        data:[...data.slice(0,index), editItem, ...data.slice(index+1)],
    }
}

const savemachineData = (oldState, save) =>{
    const {data} = oldState
    let saveItem = {
        ...save.editInput,
        key: data[save.index].key,
        edit:false
    }
    return {
        ...oldState,
        data: [...data.slice(0,save.index), saveItem, ...data.slice(save.index+1)],
        editState: false,
    }
}

const delmachineData =(oldState, index) =>{
    const {data, page} = oldState
    let add =0
    if(index%page.pageSize==0 && page.current > 1){
      add = -1
    }
    return {
        ...oldState,
        data: [...data.slice(0,index), ...data.slice(index+1)],
        page:{
          pageSize: page.pageSize,
          current:  page.current+add
        }
     }
}

const machines = (state= machinesInitState, action)=>{
    switch(action.type){
        case machine.REQ_MACHINES:
            return state;
        case machine.RECV_MACHINES:
            return {
                ...state,
                data: action.data,
                cur: action.data.length,
            }
        case machine.ADD_MACHINE:
            return addmachineData(state, action.newItem)
        case machine.EDIT_MACHINE:
            return editmachineData(state, action.index)
        case machine.SAVE_MACHINE:
            return savemachineData(state, action.save)
        case machine.DEL_MACHINE:
            return delmachineData(state, action.index)
        case machine.PAGE_MACHINE:
            return {
                ...state,
                page:action.page,
            }
        default:
            return state;
    }
}

const rootReducer = combineReducers({
    layout,
    machines,
})

export default rootReducer