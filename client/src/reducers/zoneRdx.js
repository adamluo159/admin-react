import {zoneActions} from '../actions/zoneAction'

const zoneInitState = {
    zoneEdit: true,
}

const objReduxHandle = (oldState, playload) => {
    return {
        ...oldState,
        ...playload
    }
}

const zoneRdxHandle = {}
zoneRdxHandle[zoneActions.DISABLE_ZONE_EDIT] = objReduxHandle
//zoneRdxHandle[zoneActions.ADD_ZONE] = addmachineData
//zoneRdxHandle[zoneActions.EDIT_ZONE] = editmachineData
//zoneRdxHandle[zoneActions.DEL_ZONE] = delmachineData
//zoneRdxHandle[zoneActions.SAVE_ZONE] = savemachineData

const zone = (state = zoneInitState, action) => {
    let handle = zoneRdxHandle[action.type]
    if (handle) {
        return handle(state, action.playload)
    } else {
        return state
    }
}

export default zone
