import {combineReducers} from 'redux'
import layout from './layoutRdx'
import machines from './machineRdx'
import zone from './zoneRdx'

const rootReducer = combineReducers({layout, machines, zone})
export default rootReducer