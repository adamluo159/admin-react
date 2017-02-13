import {combineReducers} from 'redux'
import layout from './layoutRdx'
import machines from './machineRdx'

const rootReducer = combineReducers({layout, machines})
export default rootReducer