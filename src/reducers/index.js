import { combineReducers } from 'redux'
import {REQ_MACHINES, 
        RECV_MACHINES,
        ADD_MACHINE,
        EDIT_MACHINE,
        SAVE_MACHINE,
        SELECT_MAINLAYOUT_KEY,
} from '../actions'

const initState = {
    loading : true,
}
const layout = (state ={selectKey: 'machineMgr'}, action)=>{
    switch (action.type) {
        case SELECT_MAINLAYOUT_KEY:
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

const machines = (state= machinesInitState, action)=>{
    switch(action.type){
        case REQ_MACHINES:
            return state;
        case RECV_MACHINES:
            return {
                ...state,
                data: action.data,
                cur: action.data.length,
            }
        case ADD_MACHINE:
            return addmachineData(state, action.newItem)
        case EDIT_MACHINE:
            return editmachineData(state, action.index)
        case SAVE_MACHINE:
            return savemachineData(state, action.save)
        default:
            return state;
    }
}

const rootReducer = combineReducers({
    layout,
    machines,
})

export default rootReducer