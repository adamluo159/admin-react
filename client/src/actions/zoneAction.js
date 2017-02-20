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
            .then(json => playload.addZone(json))
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
            .then(json => playload.saveZone({ oldzid: playload.oldZid, json: json }))
    }
}
const fetchSynMachine = (obj) => {
    console.log("syn:", obj)
    return dispatch => {
        return fetch("/zone/synMachine?zid=" + obj.zid, {
            method: "GET",
        })
            .then(response =>response.json())
            .then(json => console.log("synRsp:", json))
    }

}



export const zoneActions = {}
const mapZone = {
    //网络请求的action
    "fetchInitZones": fetchInitZones,
    "fetchAddZone": fetchAddZone,
    "fetchSaveZone": fetchSaveZone,
    "fetchSynMachine": fetchSynMachine,
    //"fetchDelZone": fetchDelZone
}
export default mapZone