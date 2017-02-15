import fetch from 'isomorphic-fetch'
import {create} from '../utils/utils'

import  machine from './machineAction'
import  layout  from './layoutAction'
import  zone    from './zoneAction'

const dispatchFunc = dispatch => {
    let mapDispatch = {}
    mapDispatch["machineD"] = create(machine, dispatch)
    mapDispatch["layoutD"] = create(layout, dispatch)
    mapDispatch["zoneD"] = create(zone, dispatch)
    return mapDispatch
}

export default dispatchFunc