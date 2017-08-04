import api from '../api'

export const INIT_MACHINES = 'INIT_MACHINES';
export const ADD_MACHINE = `ADD_MACHINE`;
export const EDIT_MACHINE = `EDIT_MACHINE`;
export const FILTER_MACHINES = `FILTER_MACHINES`;

export const fetchInitMachines = () => {
    return {
            promise: api.get('/machine')
    }
}

export const fetchSaveMachine = (playload) => {
    //    return dispatch => {
    //        let body = JSON.stringify({
    //            oldhost: playload.oldmachine.hostname,
    //            Item: playload.machine
    //        })
    //        return fetch("/machine/save", {
    //            method: "POST",
    //            headers: {
    //                "Content-Type": "application/json"
    //            },
    //            body,
    //        })
    //            .then(response => response.json())
    //            .then(json => playload.cb(json))
    //    }
    //
    return {
        type: 'MACHINE',
        payload: {
            promise: api.post('/machine/save', playload)
        }
    }
}

export const fetchAddMachine = (playload) => {
    //return dispatch => {
    //    //dispatch(machineDispatch.)
    //    return fetch("/machine/add", {
    //        method: "POST",
    //        headers: {
    //            "Content-Type": "application/json"
    //        },
    //        body: JSON.stringify(playload.machine)
    //    })
    //        .then(response => response.json())
    //        .then(json => rspAddMachine(dispatch, playload, json))
    //}
    return {
        type: 'MACHINE',
        payload: {
            promise: api.post('/machine/add', playload)
        }
    }
}

export const fetchDelMachine = (playload) => {
    //return dispatch => {
    //    return fetch("/machine/del", {
    //        method: "POST",
    //        headers: {
    //            "Content-Type": "application/json"
    //        },
    //        body: JSON.stringify(playload.fetchDel)
    //    })
    //        .then(response => response.json())
    //        .then(json => playload.delCB(json))
    //}
    return {
        type: 'MACHINE',
        payload: {
            promise: api.post('/machine/del', playload)
        }
    }
}

export const fetchCommonConfig = (f) => {
    //    return dispatch => {
    //        return fetch("/machine/common", )
    //            .then(response => response.json())
    //            .then(json => f(json))
    //    }
    return {
        type: 'MACHINE',
        payload: {
            promise: api.post('/machine/common', playload)
        }
    }

}

export const fetchSvnUpdate = (playload) => {
    //return dispatch => {
    //    return fetch("/machine/svnUpdate", {
    //        method: "POST",
    //        headers: {
    //            "Content-Type": "application/json"
    //        },
    //        body: JSON.stringify(playload.obj)
    //    })
    //        .then(response => response.json())
    //        .then(json => playload.cb(json))
    //}
    return {
        type: 'MACHINE',
        payload: {
            promise: api.post('/machine/svnUpdate', playload)
        }
    }
}

export const fetchAllMachineSvnUpdate = (f) => {
    //return dispatch => {
    //    return fetch("/machine/svnUpdateAll", )
    //        .then(response => response.json())
    //        .then(json => f(json))
    //}
    return {
        type: 'MACHINE',
        payload: {
            promise: api.post('/machine/svnUpdateAll', playload)
        }
    }
}