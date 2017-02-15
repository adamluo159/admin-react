import {zoneActions} from '../actions/zoneAction'

const zoneInitState = {
}

const zoneInitData = (oldState, init) =>{
    return {
        init
    }
}


const zoneRdxHandle = {}
//zoneRdxHandle[zoneActions.INIT_ZONES] = objReduxHandle
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
