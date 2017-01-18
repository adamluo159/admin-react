import {
	combineReducers
} from 'redux'
import {
	actions
} from '../actions'

const initState = {
	loading: true
}
const layout = (state = {
	selectKey: 'machineMgr'
}, action) => {
	switch (action.type) {
		case actions.SELECT_MAINLAYOUT_KEY:
			return {
				selectKey: action.playload
			}
		default:
			return state
	}
}

const machinesInitState = {
	editState: false,
	editInput: {},
	data: [],
	page: {
		current: 1,
		pageSize: 3
	},
	cur: 0
}

const addmachineData = (oldState, newItem) => {
	const {
		cur,
		page,
		data
	} = oldState
	newItem.key += cur
	newItem.hostname += cur
	return {
		editState: true,
		page: page,
		cur: cur + 1,
		data: [
			...data,
			newItem
		],
		editInput: newItem
	}
}

const editmachineData = (oldState, index) => {
	const {
		data
	} = oldState
	let editItem = Object.assign({}, data[index])
	editItem.edit = true
	return {
		...oldState,
		editState: true,
		editInput: {
			...oldState[index]
		},
		data: [
			...data.slice(0, index),
			editItem,
			...data.slice(index + 1)
		]
	}
}

const savemachineData = (oldState, save) => {
	const {
		data
	} = oldState
	let saveItem = {
		...save.editInput,
		key: data[save.index].key,
		edit: false
	}
	return {
		...oldState,
		data: [
			...data.slice(0, save.index),
			saveItem,
			...data.slice(save.index + 1)
		],
		editState: false
	}
}

const delmachineData = (oldState, index) => {
	const {
		data,
		page
	} = oldState
	let add = 0
	if (index % page.pageSize == 0 && page.current > 1) {
		add = -1
	}
	return {
		...oldState,
		data: [
			...data.slice(0, index),
			...data.slice(index + 1)
		],
		page: {
			pageSize: page.pageSize,
			current: page.current + add
		}
	}
}

const objReduxHandle = (state, playload) => ({
	...state,
	...playload
})
const reduxHandle = {}
reduxHandle[actions.RECV_MACHINES] = objReduxHandle
reduxHandle[actions.ADD_MACHINE] = addmachineData
reduxHandle[actions.EDIT_MACHINE] = editmachineData
reduxHandle[actions.PAGE_MACHINE] = objReduxHandle
reduxHandle[actions.DEL_MACHINE] = delmachineData
reduxHandle[actions.SAVE_MACHINE] = savemachineData

const machines = (state = machinesInitState, action) => {
	let handle = reduxHandle[action.type]
	console.log("xswxsw", action)
	if (handle) {
		return handle(state, action.playload)
	} else {
		return state
	}
}
const rootReducer = combineReducers({
	layout,
	machines
})

export default rootReducer