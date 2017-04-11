import fetch from 'isomorphic-fetch'
import {actionCreator} from '../utils/utils'

const rspInitMachines =(dispatch, initFunc, rsp) =>{
    if (rsp.Items== null){
        return 
    }
    //rsp.Items.forEach(element =>{
    //    element.applications = element.applications.toString()
    //});
    dispatch(mapMachine.InitMachines({
        data: rsp.Items,
        editState: false
    }))
    initFunc()
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
    if (rsp.Result === "OK") {
        playload.cb(rsp.Item.hostname,rsp.oldhost)
        dispatch(mapMachine.saveMachine({index:playload.index, rsp}))
    } else {
        let rsp = {
            Item:playload.oldmachine
        }
        rsp.Item.edit = false
        dispatch(mapMachine.saveMachine({index:playload.index, rsp}))
    }
}

const fetchSaveMachine = (playload) => {
    return dispatch => {
        let body = JSON.stringify({
            oldhost:playload.oldmachine.hostname,
            Item: playload.machine
        })
        return fetch("/machine/save", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body,
            })
            .then(response => response.json())
            .then(json=> rspSaveMachine(dispatch, playload, json))
    }
}

const rspAddMachine = (dispatch, playload, rsp) => {
    if (rsp.Result === "OK") {
        playload.cb(rsp.Item.hostname)
        dispatch(mapMachine.saveMachine({index:playload.index, rsp}))
    } else {
        dispatch(mapMachine.delMachine(playload.index))
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
        .then(json => rspAddMachine(dispatch, playload, json))
    }
}

const rspDelMachine = (dispatch, playload, rsp) =>{
    if(rsp.Result === "OK"){
        playload.delCB()
        dispatch(mapMachine.delMachine(playload.index))
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
        .then(json=>rspDelMachine(dispatch, playload, json))
    }
}

export const machineActions = {}
const mapMachine = {
    //界面表现的action
    "reqMachines": actionCreator(machineActions, 'REQ_MACHINES'),
    "InitMachines": actionCreator(machineActions, 'INIT_MACHINES'),
    "addMachine": actionCreator(machineActions, 'ADD_MACHINE'),
    "editMachine": actionCreator(machineActions, 'EDIT_MACHINE'),
    "saveMachine": actionCreator(machineActions, 'SAVE_MACHINE'),
    "delMachine": actionCreator(machineActions, 'DEL_MACHINE'),
    "pageMachine": actionCreator(machineActions, 'PAGE_MACHINE'),

    //网络请求的action
    "fetchInitMachines": fetchInitMachines,
    "fetchSaveMachine": fetchSaveMachine,
    "fetchAddMachine": fetchAddMachine,
    "fetchDelMachine": fetchDelMachine,
}

export default mapMachine;

