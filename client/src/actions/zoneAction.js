import fetch from 'isomorphic-fetch'
import { actionCreator } from '../utils/utils'

const fetchInitZones = (initFunc) => {
    return dispatch => {
        //dispatch(machineDispatch.reqZones())
        return fetch("/zone", )
            .then(response => response.json())
            .then(json => initFunc(json))
    }
}

const fetchAddZone = (playload) => {
    return dispatch => {
        //dispatch(machineDispatch.reqZones())
        return fetch("/zone/add", {
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
const fetchSaveZone = (playload) => {
    return dispatch => {
        //dispatch(machineDispatch.reqZones())
        let body = JSON.stringify({
            OldZoneName: playload.oldZoneName,
            OldZid: playload.oldZid,
            Item: playload.obj
        })
        return fetch("/zone/save", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body
        })
            .then(response => response.json())
            .then(json => playload.cb({ oldzid: playload.oldZid, json: json }))
    }
}
const fetchSynMachine = (obj) => {
    return dispatch => {
        return fetch("/zone/synMachine?zid=" + obj.zid + "&hostname=" + obj.hostname, {
            method: "GET",
        })
            .then(response => response.json())
            .then(json => obj.cb(json))

    }
}
const fetchDelZone = (playload) => {
    return dispatch => {
        return fetch("/zone/del", {
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

const fetchUpdateZonelogdb = (playload) => {
    console.log("updateZonelogdb, ", playload)
    return dispatch => {
        //dispatch(machineDispatch.reqZones())
        return fetch("/zone/updateZonelogdb", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(playload)
        })
            .then(response => response.json())
            .then(json => console.log("updateZonelogdb:", json))
    }
}

const fetchStartZone = (playload) => {
    return dispatch => {
        console.log("fetch start:", playload)
        return fetch("/zone/startZone", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(playload.obj)
        })
            .then(response => response.json())
            .then(json => playload.startZoneRsp(json))
    }
}

const fetchStopZone = (playload) => {
    return dispatch => {
        console.log("fetch stop:", playload)
        return fetch("/zone/stopZone", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(playload.obj)
        })
            .then(response => response.json())
            .then(json => playload.stopZoneRsp(json))
    }
}

const fetchStartAllZone = (cb) => {
    return dispatch => {
        return fetch("/zone/startAllZone", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
        })
            .then(response => response.json())
            .then(json => cb(json))
    }
}
const fetchStopAllZone = (cb) => {
    return dispatch => {
        return fetch("/zone/stopAllZone", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
        })
            .then(response => response.json())
            .then(json => cb(json))
    }
}


export const zoneActions = {}
const mapZone = {
    //网络请求的action
    "fetchInitZones": fetchInitZones,
    "fetchAddZone": fetchAddZone,
    "fetchSaveZone": fetchSaveZone,
    "fetchSynMachine": fetchSynMachine,
    "fetchStartZone": fetchStartZone,
    "fetchStopZone": fetchStopZone,
    "fetchDelZone": fetchDelZone,
    "fetchUpdateZonelogdb": fetchUpdateZonelogdb,
    "fetchStartAllZone": fetchStartAllZone,
    "fetchStopAllZone": fetchStopAllZone,
    "DisableEdit": actionCreator(zoneActions, 'DISABLE_ZONE_EDIT'),
}
export default mapZone