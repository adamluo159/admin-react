import { combineReducers } from 'redux'
import {REQ_LAYOUT_POSTS, RECV_LAYOUT_POSTS} from '../actions'

const initState = {
    loading : true,
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
    fetchLayoutData
})

export default rootReducer