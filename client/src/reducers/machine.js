import {
    INIT_MACHINES,
    ADD_MACHINE,
    EDIT_MACHINE,
    FILTER_MACHINES
} from '../actions/machine'

const InitState = {
    //data: [],
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
machineRdxHandle[INIT_MACHINES] = objReduxHandle
machineRdxHandle[ADD_MACHINE] = addmachineData
machineRdxHandle[EDIT_MACHINE] = editmachineData
machineRdxHandle[FILTER_MACHINES] = filterMachine

const machines = (state = InitState, action) => {
    let handle = machineRdxHandle[action.type]
    if (handle) {
        return handle(state, action.playload)
    } else {
        return state
    }
}
export default machines
