import {machineActions} from '../actions/machineAction'

const machinesInitState = {
    data: [],
    page: {
        current: 1,
        pageSize: 30
    }
}

const addmachineData = (oldState, newItem) => {
    const {page, data} = oldState
    return {
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

const objReduxHandle = (oldState, playload) => {
    return {
        ...oldState,
        ...playload
    }
}

const filterMachine = (oldState, data) => {
    return {
        ...oldState,
        data:data,
    }
}

const machineRdxHandle = {}
machineRdxHandle[machineActions.INIT_MACHINES] = objReduxHandle
machineRdxHandle[machineActions.ADD_MACHINE] = addmachineData
machineRdxHandle[machineActions.EDIT_MACHINE] = editmachineData
machineRdxHandle[machineActions.PAGE_MACHINE] = objReduxHandle
machineRdxHandle[machineActions.FILTER_MACHINES] = filterMachine

const machines = (state = machinesInitState, action) => {
    let handle = machineRdxHandle[action.type]
    if (handle) {
        return handle(state, action.playload)
    } else {
        return state
    }
}
export default machines
