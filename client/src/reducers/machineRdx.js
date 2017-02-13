import {machineActions} from '../actions/machineAction'

const machinesInitState = {
    editState: false,
    data: [],
    page: {
        current: 1,
        pageSize: 30
    }
}

const addmachineData = (oldState, newItem) => {
    const {page, data} = oldState
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
        },
        editState: false
    }
}

const objReduxHandle = (oldState, playload) => {
    return {
        ...oldState,
        ...playload
    }
}

const machineRdxHandle = {}
machineRdxHandle[machineActions.INIT_MACHINES] = objReduxHandle
machineRdxHandle[machineActions.ADD_MACHINE] = addmachineData
machineRdxHandle[machineActions.EDIT_MACHINE] = editmachineData
machineRdxHandle[machineActions.PAGE_MACHINE] = objReduxHandle
machineRdxHandle[machineActions.DEL_MACHINE] = delmachineData
machineRdxHandle[machineActions.SAVE_MACHINE] = savemachineData

const machines = (state = machinesInitState, action) => {
    let handle = machineRdxHandle[action.type]
    if (handle) {
        return handle(state, action.playload)
    } else {
        return state
    }
}
export default machines
