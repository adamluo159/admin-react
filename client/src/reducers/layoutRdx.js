import {layoutActions} from '../actions/layoutAction'
const layout = (state = {
    selectKey: 'zone'
}, action) => {
    switch (action.type) {
        case layoutActions.SELECT_MAINLAYOUT_KEY:
            return  action.playload
        default:
            return state
    }
}
export default layout