import { combineReducers } from 'redux'
import {REQ_LAYOUT_POSTS, 
        RECV_LAYOUT_POSTS,
        SELECT_MAINLAYOUT_KEY,
} from '../actions'

const initState = {
    loading : true,
}
const layout = (state ={selectKey: 'machineMgr'}, action)=>{
    switch (action.type) {
        case SELECT_MAINLAYOUT_KEY:
            return {
                selectKey: action.selectKey,
            }
        default:
            return state
    }
}


const fetchLayoutData = (state= initState, action)=>{
    switch(action.type){
        case REQ_LAYOUT_POSTS:
        return {
            loading: action.loading
        }
        case RECV_LAYOUT_POSTS:
        return {
            loading:action.loading,
            menu:action.jsonData,
        }
        default:
        return state;
    }
}

const rootReducer = combineReducers({
    layout,
    fetchLayoutData
})

export default rootReducer