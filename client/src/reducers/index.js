import {combineReducers} from 'redux'
import {actions} from '../actions'

const initState = {
    loading: true
}
const layout = (state = {
    selectKey: 'machineMgr'
}, action) => {
    switch (action.type) {
        case actions.SELECT_MAINLAYOUT_KEY:
            return action.playload
        default:
            return state
    }
}

const machinesInitState = {
    editState: false,
    data: [],
    page: {
        current: 1,
        pageSize: 3
    }
}

const addmachineData = (oldState, newItem) => {
    const {page, data} = oldState
    newItem.key = newItem.hostname
    return {
        editState: true,
        page: page,
        data: [
            ...data,
            newItem
        ],
        editInput: newItem
    }
}

const editmachineData = (oldState, index) => {
    const {data} = oldState
    let editItem = Object.assign({}, data[index])
    editItem.edit = true
    return {
        ...oldState,
        editState: true,
        editInput: {
            ...oldState[index]
        },
        data: [
            ...data.slice(0, index),
            editItem,
            ...data.slice(index + 1)
        ]
    }
}

const savemachineData = (oldState, save) => {
    const {data} = oldState
    const {index, rsp} = save
    if (rsp.Result !== "OK") {
        return
    }
    return {
        ...oldState,
        data: [
            ...data.slice(0, index),
            rsp.Item,
            ...data.slice(index + 1)
        ],
        editState: false
    }
}

const delmachineData = (oldState, index) => {
    const {data, page} = oldState
    let add = 0
    if (index % page.pageSize == 0 && page.current > 1) {
        add = -1
    }
    return {
        ...oldState,
        data: [
            ...data.slice(0, index),
            ...data.slice(index + 1)
        ],
        page: {
            pageSize: page.pageSize,
            current: page.current + add
        }
    }
}
const initMachine = (oldState, playload) => {
    if (playload.data == null) {
        return oldState
    }
    return {
        ...oldState,
        ...playload
    }
}

const objReduxHandle = (oldState, playload) => {
    return {
        ...oldState,
        ...playload
    }
}

const reduxHandle = {}
reduxHandle[actions.RECV_MACHINES] = initMachine
reduxHandle[actions.RESET_MACHINE_STATE] = objReduxHandle
reduxHandle[actions.ADD_MACHINE] = addmachineData
reduxHandle[actions.EDIT_MACHINE] = editmachineData
reduxHandle[actions.PAGE_MACHINE] = objReduxHandle
reduxHandle[actions.DEL_MACHINE] = delmachineData
reduxHandle[actions.SAVE_MACHINE] = savemachineData

const machines = (state = machinesInitState, action) => {
    let handle = reduxHandle[action.type]
    if (handle) {
        return handle(state, action.playload)
    } else {
        return state
    }
}
const rootReducer = combineReducers({layout, machines})

export default rootReducer