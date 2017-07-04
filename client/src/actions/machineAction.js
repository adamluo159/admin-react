import fetch from 'isomorphic-fetch'
import { actionCreator } from '../utils/utils'

const fetchInitMachines = (initFunc) => {
    return dispatch => {
        //dispatch(machineDispatch.reqMachines())
        return fetch("/machine", )
            .then(response => response.json())
            .then(json => initFunc(json))
    }
}

const fetchSaveMachine = (playload) => {
    return dispatch => {
        let body = JSON.stringify({
            oldhost: playload.oldmachine.hostname,
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
            .then(json => playload.cb(json))
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

const fetchDelMachine = (playload) => {
    return dispatch => {
        return fetch("/machine/del", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(playload.fetchDel)
        })
            .then(response => response.json())
            .then(json => playload.delCB(json))
    }
}

const fetchCommonConfig = (f) => {
    return dispatch => {
        return fetch("/machine/common", )
            .then(response => response.json())
            .then(json => f(json))
    }
}

const fetchSvnUpdate = (playload) => {
    return dispatch => {
        return fetch("/machine/svnUpdate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(playload.obj)
        })
            .then(response => response.json())
            .then(json => playload.cb(json))
    }
}

const fetchAllMachineSvnUpdate = (f) => {
    return dispatch => {
        return fetch("/machine/svnUpdateAll", )
            .then(response => response.json())
            .then(json => f(json))
    }
}


export const machineActions = {}
const mapMachine = {
    //界面表现的action
    "InitMachines": actionCreator(machineActions, 'INIT_MACHINES'),
    "addMachine": actionCreator(machineActions, 'ADD_MACHINE'),
    "editMachine": actionCreator(machineActions, 'EDIT_MACHINE'),
    "pageMachine": actionCreator(machineActions, 'PAGE_MACHINE'),
    "filterMachine": actionCreator(machineActions, 'FILTER_MACHINES'),

    //网络请求的action
    "fetchInitMachines": fetchInitMachines,
    "fetchSaveMachine": fetchSaveMachine,
    "fetchAddMachine": fetchAddMachine,
    "fetchDelMachine": fetchDelMachine,
    "fetchCommonConfig": fetchCommonConfig,
    "fetchSvnUpdate": fetchSvnUpdate,
    "fetchAllMachineSvnUpdate": fetchAllMachineSvnUpdate,
}


export default mapMachine;

