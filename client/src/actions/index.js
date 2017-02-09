import fetch from 'isomorphic-fetch'
const actionCreator = (regActionType, type) => {
    regActionType[type] = type
    return (playload) => ({type, playload})
}

const rspInitMachines =(dispatch, initFunc, rsp) =>{
    console.log("initmachine:", rsp)
    if (rsp.Result === "OK") {
        dispatch(machineDispatch.InitMachines({
            data: rsp.Items,
            editState: false
        }))
        initFunc()
    }
}

const fetchInitMachines = (initFunc) => {
    return dispatch => {
        //dispatch(machineDispatch.reqMachines())
        return fetch("/machine",)
            .then(response => response.json())
            .then(json => rspInitMachines(dispatch, initFunc, json))
    }
}

const rspSaveMachine = (dispatch, playload, rsp) =>{
    console.log("recv save machine:", playload, rsp)
    if (rsp.Result === "OK") {
        dispatch(machineDispatch.saveMachine({index:playload.index, rsp}))
    } else {
        let rsp = {
            Item:playload.oldmachine
        }
        rsp.Item.edit = false
        dispatch(machineDispatch.saveMachine({index:playload.index, rsp}))
    }
}
const fetchSaveMachine = (playload) => {
    return dispatch => {
        console.log("fetchSaveMachine:",playload)
        return fetch("/machine/save", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
                body: JSON.stringify(playload.machine)
            })
            .then(response => response.json())
            .then(json=> rspSaveMachine(dispatch, playload, json))
    }
}

const rspAddMachine = (dispatch, index, rsp) => {
    console.log("recv rsp AddMachine:", index, rsp)
    if (rsp.Result === "OK") {
        dispatch(machineDispatch.saveMachine({index, rsp}))
    } else {
        dispatch(machineDispatch.delMachine(index))
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
        .then(json => rspAddMachine(dispatch, playload.index, json))
    }
}

const rspDelMachine = (dispatch, index, rsp) =>{
    console.log("rspDelMachine:", rsp)
    if(rsp.Result === "OK"){
        dispatch(machineDispatch.delMachine(index))
    }
}

const fetchDelMachine =(playload)=>{
    return dispatch =>{
        return fetch("/machine/del", {
            method:"POST", 
            headers:{
                "Content-Type": "application/json"
            },
            body: JSON.stringify(playload.fetchDel)
        })
        .then(response => response.json())
        .then(json=>rspDelMachine(dispatch, playload.index, json))
    }
}

export const actions = {}
const machineDispatch = {
    //界面表现的action
    "selectMainlayout": actionCreator(actions, 'SELECT_MAINLAYOUT_KEY'),
    "reqMachines": actionCreator(actions, 'REQ_MACHINES'),
    "InitMachines": actionCreator(actions, 'INIT_MACHINES'),
    "addMachine": actionCreator(actions, 'ADD_MACHINE'),
    "editMachine": actionCreator(actions, 'EDIT_MACHINE'),
    "saveMachine": actionCreator(actions, 'SAVE_MACHINE'),
    "delMachine": actionCreator(actions, 'DEL_MACHINE'),
    "pageMachine": actionCreator(actions, 'PAGE_MACHINE'),

    //网络请求的action
    "fetchInitMachines": fetchInitMachines,
    "fetchSaveMachine": fetchSaveMachine,
    "fetchAddMachine": fetchAddMachine,
    "fetchDelMachine": fetchDelMachine,
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
