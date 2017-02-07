import fetch from 'isomorphic-fetch'
const actionCreator = (regActionType, type) => {
    regActionType[type] = type
    return (playload) => ({type, playload})
}

const fetchMachines = () => {
    return dispath => {
        dispath(machineDispatch.reqMachines())
        return fetch("/machine",)
            .then(response => response.json())
            .then(json => dispath(machineDispatch.recvMachines({data: json})))
    }
}
const fetchSaveMachine = (playload) => {
    return dispatch => {
        //dispatch(machineDispatch.)
        return fetch("/machine/save", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
                body: JSON.stringify(playload.machine)
            })
            .then(response => response.json())
            .then(json => dispatch(machineDispatch.saveMachine()))
    }
}

const fetchAddMachine = (playload) => {
    return dispatch => {
        //dispatch(machineDispatch.)
        return fetch("/machine/add", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
                body: JSON.stringify(playload.machine)
            })
            .then(response => response.json())
            .then(json => dispatch(machineDispatch.saveMachine({
                index:playload.index,
                rsp:json
            })))
    }
}


export const actions = {}
const machineDispatch = {
    //界面表现的action
    "selectMainlayout": actionCreator(actions, 'SELECT_MAINLAYOUT_KEY'),
    "reqMachines": actionCreator(actions, 'REQ_MACHINES'),
    "recvMachines": actionCreator(actions, 'RECV_MACHINES'),
    "addMachine": actionCreator(actions, 'ADD_MACHINE'),
    "editMachine": actionCreator(actions, 'EDIT_MACHINE'),
    "saveMachine": actionCreator(actions, 'SAVE_MACHINE'),
    "delMachine": actionCreator(actions, 'DEL_MACHINE'),
    "pageMachine": actionCreator(actions, 'PAGE_MACHINE'),
    "resetMachineState": actionCreator(actions, 'RESET_MACHINE_STATE'),

    //网络请求的action
    "fetchMachines": fetchMachines,
    "fetchSaveMachine": fetchSaveMachine,
    "fetchAddMachine": fetchAddMachine
}

const dispatchFunc = dispatch => {
    let dispatchObj = {}
    Object
        .keys(machineDispatch)
        .forEach(item => {
            dispatchObj[item] = (e) => dispatch(machineDispatch[item](e))
        })
    return {machineDispatch: dispatchObj}
}

export default dispatchFunc;
