import fetch from 'isomorphic-fetch'
const actionCreator = (regActionType, type) => {
    regActionType[type] = type
    return (playload) => {console.log("action:",type,playload); return {type, playload}}
}

const fetchMachines = () => {
    return dispath => {
        //dispath(machineDispatch.reqMachines())
        let data = [
            {
                key: "host0",
                hostname: 'host0',
                IP: "192.168.1.1",
                outIP: "192.168.1.1",
                type: "login",
                edit: false
            }
        ]
        dispath(machineDispatch.recvMachines({data: data}))
        // return fetch('') .then(response=>response.json())
        // .then(json=>dispath(RecvLayoutData(loading, json)))
    }
}

export const actions = {}
const machineDispatch = {
        "selectMainlayout": actionCreator(actions, 'SELECT_MAINLAYOUT_KEY'),
        "reqMachines": actionCreator(actions, 'REQ_MACHINES'),
        "recvMachines": actionCreator(actions, 'RECV_MACHINES'),
        "addMachine": actionCreator(actions, 'ADD_MACHINE'),
        "editMachine": actionCreator(actions, 'EDIT_MACHINE'),
        "saveMachine": actionCreator(actions, 'SAVE_MACHINE'),
        "delMachine": actionCreator(actions, 'DEL_MACHINE'),
        "pageMachine": actionCreator(actions, 'PAGE_MACHINE'),
        "fetchMachines": fetchMachines
}

const dispatchFunc = dispatch => {
    let dispatchObj = {}
    Object.keys(machineDispatch).forEach(item => {
    dispatchObj[item] = (e)=>dispatch(machineDispatch[item](e))})
    return {
        machineDispatch:dispatchObj
    }
}

export default dispatchFunc;
