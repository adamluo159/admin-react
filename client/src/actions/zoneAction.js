import fetch from 'isomorphic-fetch'
import {actionCreator} from '../utils/utils'

export const zoneActions = {}
const mapZone = {
    //界面表现的action
    "InitZones": actionCreator(zoneActions, 'INIT_ZONES'),
    "addZone": actionCreator(zoneActions, 'ADD_ZONE'),
    "editZone": actionCreator(zoneActions, 'EDIT_ZONE'),
    "saveZone": actionCreator(zoneActions, 'SAVE_ZONE'),
    "delZone": actionCreator(zoneActions, 'DEL_MACHINE'),

    //网络请求的action
//    "fetchInitZones": fetchInitMachines,
//    "fetchSaveZone": fetchSaveMachine,
//    "fetchAddZone": fetchAddMachine,
//    "fetchDelZone": fetchDelMachine
}
export default mapZone